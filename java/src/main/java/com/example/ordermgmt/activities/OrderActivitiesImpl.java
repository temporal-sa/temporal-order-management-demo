package com.example.ordermgmt.activities;

import com.example.ordermgmt.model.OrderInput;
import com.example.ordermgmt.model.OrderItem;
import io.temporal.activity.Activity;
import io.temporal.failure.ApplicationFailure;
import io.temporal.spring.boot.ActivityImpl;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

import java.util.List;
import java.util.concurrent.TimeUnit;

@Slf4j
@Component
@ActivityImpl(taskQueues = "${ordermgmt.task-queue}")
public class OrderActivitiesImpl implements OrderActivities {
    private static final String ERROR_CHARGE_API_UNAVAILABLE = "OrderWorkflowAPIFailure";
    private static final String ERROR_INVALID_CREDIT_CARD = "OrderWorkflowNonRecoverableFailure";

    private static void simulateExternalOperation(long ms) {
        try {
            TimeUnit.MILLISECONDS.sleep(ms);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }

    private static String simulateExternalOperation(long ms, String type, int attempt) {
        simulateExternalOperation(ms / attempt);
        return (attempt < 5) ? type : "NoError";
    }

    @Override
    public List<OrderItem> getItems() {
        log.info("Getting list of items");

        // simulate DB query
        simulateExternalOperation(100);

        return List.of(
                new OrderItem(654300, "Table Top", 1),
                new OrderItem(654321, "Table Legs", 2),
                new OrderItem(654322, "Keypad", 1)
        );
    }

    @Override
    public String checkFraud(OrderInput input) {
        log.info("Check Fraud activity started, orderId = {}", input.getOrderId());

        // simulate external API call
        simulateExternalOperation(1000);

        return input.getOrderId();
    }

    @Override
    public String prepareShipment(OrderInput input) {
        log.info("Prepare Shipment activity started, orderId = {}", input.getOrderId());

        // simulate external API call
        simulateExternalOperation(1000);

        return input.getOrderId();
    }

    @Override
    public String chargeCustomer(OrderInput input, String type) {
        log.info("Charge Customer activity started, orderId = {}", input.getOrderId());
        int attempt = Activity.getExecutionContext().getInfo().getAttempt();

        // simulate external API call
        String error = simulateExternalOperation(1000, type, attempt);
        switch (error) {
            case ERROR_CHARGE_API_UNAVAILABLE:
                // a transient error, which can be retried
                log.info("Charge Customer API unavailable, attempt = {}", attempt);
                throw new RuntimeException("Charge Customer activity failed, API unavailable");
            case ERROR_INVALID_CREDIT_CARD:
                // a business error, which cannot be retried
                throw ApplicationFailure.newNonRetryableFailure("Charge Customer activity failed, card is invalid", "InvalidCreditCard");
            default:
                // pass through, no error
        }

        return input.getOrderId();
    }

    @Override
    public void shipOrder(OrderInput input, OrderItem item) {
        log.info("Ship Order activity started, orderId ={}, itemId = {}, itemDescription = {}", input.getOrderId(), item.getId(), item.getDescription());

        // simulate external API call
        simulateExternalOperation(1000);
    }

    @Override
    public String undoPrepareShipment(OrderInput input) {
        log.info("Undo Prepare Shipment activity started, orderId = {}", input.getOrderId());

        // simulate external API call
        simulateExternalOperation(1000);

        return input.getOrderId();
    }

    @Override
    public String undoChargeCustomer(OrderInput input) {
        log.info("Undo Charge Customer activity started, orderId = {}", input.getOrderId());

        // simulate external API call
        simulateExternalOperation(1000);

        return input.getOrderId();
    }
}
