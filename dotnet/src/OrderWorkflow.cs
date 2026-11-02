using System.Collections;
using DotNetOrderManagement.model;
using Microsoft.Extensions.Logging;
using Temporalio.Workflows;

namespace DotNetOrderManagement;

[Workflow("OrderWorkflowHappyPath")]
public class OrderWorkflow
{
    private int _progress = 0;

    [WorkflowRun]
    public async Task<OrderOutput> Execute(OrderInput input)
    {
        var logger = Workflow.Logger;
        var type = Workflow.Info.WorkflowType;
        logger.LogInformation("Order workflow started, type = {}, orderId = {}", type, input.OrderId);

        // Get Items
        var orderItems = await Workflow.ExecuteLocalActivityAsync((OrderActivities act) =>
            act.GetItems(), OrderActivities.LocalActivityOpts);

        // Check fraud
        await Workflow.ExecuteActivityAsync((OrderActivities act) =>
            act.CheckFraud(input), OrderActivities.ActivityOpts);
        await Sleep(1, 25);

        // Prepare shipment
        await Workflow.ExecuteActivityAsync((OrderActivities act) =>
            act.PrepareShipment(input), OrderActivities.ActivityOpts);
        await Sleep(1, 50);

        // Charge Customer
        await Workflow.ExecuteActivityAsync((OrderActivities act) =>
            act.ChargeCustomer(input, type), OrderActivities.ActivityOpts);
        await Sleep(3, 75);

        // Ship order items
        var shipOpts = (ActivityOptions)OrderActivities.ActivityOpts.Clone();
        // item n takes n*2s; leave headroom so start-to-close doesn't fire first
        shipOpts.StartToCloseTimeout = TimeSpan.FromSeconds(15);
        // a dead worker is detected in ~2s instead of waiting out the full start-to-close
        shipOpts.HeartbeatTimeout = TimeSpan.FromSeconds(2);

        var activities = new List<Task>();
        for (var i = 0; i < orderItems.Count; i++)
        {
            var shippingInput = new ShippingInput(input, orderItems[i], i + 1);
            logger.LogInformation("Shipping item: {}", shippingInput.Item.Description);
            activities.Add(
                Workflow.ExecuteActivityAsync((OrderActivities act) =>
                    act.ShipOrder(shippingInput), shipOpts)
            );
        }

        // Wait for all items to ship
        await Workflow.WhenAllAsync(activities);
        await Sleep(0, 100);

        // Generate tracking ID
        var trackingId = Workflow.Random.Next().ToString();
        return new OrderOutput(trackingId, input.Address);
    }

    [WorkflowQuery("getProgress")]
    public int QueryProgress()
    {
        return _progress;
    }

    private async Task Sleep(int sleep, int progress)
    {
        _progress = progress;
        if (sleep > 0)
        {
            await Workflow.DelayAsync(TimeSpan.FromSeconds(sleep));
        }
    }
}
