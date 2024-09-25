using System.Diagnostics;
using DotNetOrderManagement.model;
using Microsoft.Extensions.Logging;
using System.Threading.Tasks;

using Temporalio.Common;
using Temporalio.Converters;
using Temporalio.Exceptions;
using Temporalio.Workflows;

namespace DotNetOrderManagement;

[Workflow(Dynamic = true)]
public class OrderWorkflowScenarios
{
    private static readonly string BUG = "OrderWorkflowRecoverableFailure";
    private static readonly string CHILD = "OrderWorkflowChildWorkflow";
    private static readonly string SIGNAL = "OrderWorkflowHumanInLoopSignal";
    private static readonly string UPDATE = "OrderWorkflowHumanInLoopUpdate";
    private static readonly string VISIBILITY = "OrderWorkflowAdvancedVisibility";
    
    private SearchAttributeKey<string> ORDER_STATUS_SA = SearchAttributeKey.CreateKeyword("OrderStatus");

    private int _progress = 0;
    private string _updatedAddress = string.Empty;

    [WorkflowRun]
    public async Task<OrderOutput> RunAsync(IRawValue[] args)
    {
        var logger = Workflow.Logger;
        var type = Workflow.Info.WorkflowType;
        var input = Workflow.PayloadConverter.ToValue<OrderInput>(args[0]);
        logger.LogInformation("Dynamic Order workflow started, type = {}, orderId = {}", type, input.OrderId);
        
        // Saga compensations
        var compensations = new List<Func<Task>>();

        // Get Items
        var orderItems = await Workflow.ExecuteLocalActivityAsync((OrderActivities act) => 
            act.GetItems(), OrderActivities.LocalActivityOpts);

        await UpdateProgress("Check Fraud", 0, 0);
        
        // Check fraud
        await Workflow.ExecuteActivityAsync((OrderActivities act) =>
            act.CheckFraud(input), OrderActivities.ActivityOpts);
        
        await UpdateProgress("Prepare Shipment", 25, 1);
        
        // Prepare shipment
        compensations.Add(async () => 
            await Workflow.ExecuteActivityAsync(() => 
                OrderActivities.UndoPrepareShipment(input), OrderActivities.ActivityOpts));
        
        await Workflow.ExecuteActivityAsync((OrderActivities act) => 
            act.PrepareShipment(input), OrderActivities.ActivityOpts);
        
        await UpdateProgress("Charge Customer", 50, 1);
        
        // Charge Customer
        try
        {
            compensations.Add(async () => 
                await Workflow.ExecuteActivityAsync(() => 
                    OrderActivities.UndoChargeCustomer(input), OrderActivities.ActivityOpts));
            
            await Workflow.ExecuteActivityAsync((OrderActivities act) =>
                act.ChargeCustomer(input, type), OrderActivities.ActivityOpts);
        }
        catch (ActivityFailureException af)
        {
            logger.LogError("Failed to charge customer {}", af);
            await SagaCompensate(compensations);
            throw;
        }
        
        await UpdateProgress("Ship Order", 75, 3);

        if (BUG.Equals(type))
        {
            // Simulate a bug
            // throw new Exception("Simulated bug - fix me!");
        }

        if (SIGNAL.Equals(type) || UPDATE.Equals(type))
        {
            // await message to update address
            await WaitForUpdatedAddressOrTimeout(input);
        }

        // Ship order items
        if (CHILD.Equals(type))
        {
            await ShipOrderChildWorkflowsAndWait(input, orderItems, type);
        }
        else
        {
            await ShipOrderActivitiesAndWait(input, orderItems, type);
        }
        logger.LogInformation("Items have been shipped");
        await UpdateProgress("Order Completed", 100, 2);    
        
        // Generate tracking ID
        var trackingId = Workflow.Random.Next().ToString();
        return new OrderOutput(trackingId, input.Address);
    }
    
    [WorkflowQuery("getProgress")]
    public int QueryProgress()
    {
        return _progress;
    }
    
    [WorkflowSignal("UpdateOrder")]
    public async Task UpdateOrderSignal(UpdateOrderInput updateInput)
    {
        Workflow.Logger.LogInformation("Received update order signal with address: {}", updateInput.Address);
        _updatedAddress = updateInput.Address;
    }

    [WorkflowUpdate("UpdateOrder")]
    public async Task<string> UpdateOrderUpdate(UpdateOrderInput updateInput)
    {
        Workflow.Logger.LogInformation("Received update order update with address: {}", updateInput.Address);
        _updatedAddress = updateInput.Address;
        return "Updated address: " +_updatedAddress;
    }

    [WorkflowUpdateValidator("UpdateOrderUpdate")]
    public void UpdateOrderValidator(UpdateOrderInput updateInput)
    {
        if (!Char.IsDigit(updateInput.Address[0]))
        {
            Workflow.Logger.LogInformation("Rejecting order update, invalid address: {}", updateInput.Address);
            throw new ApplicationFailureException("Address must start with a digit", "invalid-address");
        }
        
        Workflow.Logger.LogInformation("Order update address is valid: {}", updateInput.Address);
    }

    private async Task WaitForUpdatedAddressOrTimeout(OrderInput input)
    {
        Workflow.Logger.LogInformation("Waiting up to 60 seconds for updated address");
        var ok = await Workflow.WaitConditionAsync(
            () =>  !string.IsNullOrEmpty(_updatedAddress), TimeSpan.FromSeconds(60));
        if (ok)
        {
            input.Address = _updatedAddress;
        }
        else
        {
            // Do nothing - use the original address
            // In other cases, you may want to throw an exception on timeout, e.g.
            // throw new ApplicationFailureException("Updated address was not received in 60 seconds");
        }
        
    }
    
    private async Task ShipOrderActivitiesAndWait(OrderInput input, List<OrderItem> items, string type)
    {
        var logger = Workflow.Logger;
        var activities = new List<Task>();
        foreach (var orderItem in items)
        {
            logger.LogInformation("Shipping item: {}", orderItem.Description);
            activities.Add(
                Workflow.ExecuteActivityAsync((OrderActivities act) => 
                    act.ShipOrder(input, orderItem), OrderActivities.ActivityOpts)
            );
        }
        // Wait for all items to ship
        logger.LogInformation("Waiting for shipping activities to complete...");
        await Workflow.WhenAllAsync(activities);
        logger.LogInformation("Shipping activities are now completed");
    }
    
    private async Task ShipOrderChildWorkflowsAndWait(OrderInput input, List<OrderItem> items, string type)
    {
        Debug.Assert(CHILD.Equals(type));   
        var logger = Workflow.Logger;
    
        var childWorkflows = new List<Task>();
        foreach (var orderItem in items)
        {
            logger.LogInformation("Shipping item via child workflow: {}, ", orderItem.Id);
            var handle = await Workflow.StartChildWorkflowAsync(
                (ShippingChildWorkflow wf) => wf.Execute(input, orderItem), new()
                {
                    Id = string.Join("-", "shipment", input.OrderId, orderItem.Id),
                    ParentClosePolicy = ParentClosePolicy.Terminate,
                });
            
            childWorkflows.Add(handle.GetResultAsync());
        }
    
        logger.LogInformation("Waiting for child shipping workflows to complete...");
        await Workflow.WhenAllAsync(childWorkflows);
        logger.LogInformation("Shipping child workflows are now completed");
    }
    
    private async Task SagaCompensate(List<Func<Task>> compensations)
    {
        compensations.Reverse(0, compensations.Count);
        foreach (var comp in compensations)
        {
#pragma warning disable CA1031
            try
            {
                await comp.Invoke();
            }
            catch (Exception ex)
            {
                Workflow.Logger.LogError(ex, "Failed to compensate");
                // don't propagate the exception
            }
        }
#pragma warning restore CA1031
    }

    private async Task UpdateProgress(string orderStatus, int progress, int sleep)
    {
        _progress = progress;
        if (sleep > 0)
        {
            await Workflow.DelayAsync(TimeSpan.FromSeconds(sleep));
        }

        if (VISIBILITY.Equals(Workflow.Info.WorkflowType))
        {
            Workflow.Logger.LogInformation("Advanced visibility .. {}", orderStatus);
            Workflow.UpsertTypedSearchAttributes(ORDER_STATUS_SA.ValueSet(orderStatus));    
        }
    }

}