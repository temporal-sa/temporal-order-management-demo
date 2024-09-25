using System.Diagnostics;
using Microsoft.Extensions.Logging;
using Temporalio.Activities;
using Temporalio.Exceptions;

using DotNetOrderManagement.model;
using Temporalio.Workflows;

namespace DotNetOrderManagement;

public class OrderActivities
{
    private static readonly string ERROR_CHARGE_API_UNAVAILABLE = "OrderWorkflowAPIFailure";
    private static readonly string ERROR_INVALID_CREDIT_CARD = "OrderWorkflowNonRecoverableFailure";
    
    public static readonly ActivityOptions ActivityOpts = new()
    {
        StartToCloseTimeout = TimeSpan.FromSeconds(5),
        RetryPolicy = new()
        {
            InitialInterval = TimeSpan.FromSeconds(1),
            MaximumInterval = TimeSpan.FromSeconds(30),
            BackoffCoefficient = 2
        }
    };

    public static readonly LocalActivityOptions LocalActivityOpts = new()
    {
        StartToCloseTimeout = TimeSpan.FromSeconds(5),
    };    

    private static async Task<string> SimulateExternalOperationAsync(int ms)
    {
        await Task.Delay(ms);
        return "SUCCESS";
    }

    private static async Task<string> SimulateExternalOperationAsync(int ms, string type, int attempt)
    {
        _ = await SimulateExternalOperationAsync(ms / attempt);
        return (attempt < 5) ? type : "NoError";
    }

    [Activity]
    public async Task<List<OrderItem>> GetItems()
    {
        var logger = ActivityExecutionContext.Current.Logger;
        logger.LogInformation("Getting list of items");

        // simulate DB Query
        await SimulateExternalOperationAsync(100);

        return  
            [
                new OrderItem { Id = 654300, Description = "Table Top",  Quantity = 1 },
                new OrderItem { Id = 654321, Description = "Table Legs", Quantity = 2 },
                new OrderItem { Id = 654322, Description = "Keypad",     Quantity = 1 },
            ];
    }

    [Activity]
    public async Task<string> CheckFraud(OrderInput input)
    {
        var logger = ActivityExecutionContext.Current.Logger;
        logger.LogInformation("Check Fraud activity started, orderId = {}", input.OrderId);

        // simulate external operation
        await SimulateExternalOperationAsync(1000);

        return input.OrderId;
    }

    [Activity]
    public async Task<string> PrepareShipment(OrderInput input)
    {
        var logger = ActivityExecutionContext.Current.Logger;
        logger.LogInformation("Prepare Shipment activity started, orderId = {}", input.OrderId);
        
        // simulate external API call
        await SimulateExternalOperationAsync(1000);

        return input.OrderId;
    }

    [Activity]
    public async Task<string> ChargeCustomer(OrderInput input, string type)
    {
        var logger = ActivityExecutionContext.Current.Logger;
        logger.LogInformation("Charge Customer activity started, orderId = {}", input.OrderId);
        
        // simulate external API call
        var attempt = ActivityExecutionContext.Current.Info.Attempt;
        string error = await SimulateExternalOperationAsync(1000, type, attempt);
        logger.LogInformation("Simulated call complete, type {}, error = {}", type, error);
        if (ERROR_CHARGE_API_UNAVAILABLE.Equals(error))
        {
            // a transient error, which can be retried
            logger.LogInformation("Charge Customer API unavailable, attempt = {}", attempt);
            throw new Exception("Charge Customer activity failed. API unavailable.");
        }

        if (ERROR_INVALID_CREDIT_CARD.Equals(error))
        {
            // a business error, which cannot be retried
            throw new ApplicationFailureException("Charge Customer activity failed. Card is invalid",
                "InvalidCreditCard", true);
        }

        return input.OrderId;
    }

    [Activity]
    public async Task ShipOrder(OrderInput input, OrderItem item)
    {
        var logger = ActivityExecutionContext.Current.Logger;
        logger.LogInformation("Ship Order activity started, orderId ={}, itemId = {}, itemDescription = {}", 
            input.OrderId, item.Id, item.Description);
        
        // simulate external API call
        await SimulateExternalOperationAsync(1000);
    }

    [Activity]
    public static async Task<string> UndoPrepareShipment(OrderInput input)
    {
        var logger = ActivityExecutionContext.Current.Logger;
        logger.LogInformation("Undo Prepare Shipment activity started, orderId = {}", input.OrderId);
        
        // simluate external API call
        await SimulateExternalOperationAsync(1000);
        
        return input.OrderId;
    }

    [Activity]
    public static async Task<string> UndoChargeCustomer(OrderInput input)
    {
        var logger = ActivityExecutionContext.Current.Logger;
        logger.LogInformation("Undo Charge Customer activity started, orderId = {}", input.OrderId);
        
        // simulate external API call
        await SimulateExternalOperationAsync(1000);
        
        return input.OrderId;
    }

}