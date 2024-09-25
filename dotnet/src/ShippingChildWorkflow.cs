using DotNetOrderManagement.model;
using Microsoft.Extensions.Logging;
using Temporalio.Workflows;

namespace DotNetOrderManagement;

[Workflow]
public class ShippingChildWorkflow
{
    [WorkflowRun]
    public async Task Execute(OrderInput input, OrderItem orderItem)
    {
        Workflow.Logger.LogInformation("Shipping workflow started, orderId = {}", input.OrderId);
        
        // Ship order
        await Workflow.ExecuteActivityAsync((OrderActivities act) =>
            act.ShipOrder(input, orderItem), OrderActivities.ActivityOpts);
    }
}