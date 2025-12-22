package com.example.ordermgmt.workflows;

import com.example.ordermgmt.OrderApplication;
import com.example.ordermgmt.activities.OrderActivities;
import com.example.ordermgmt.model.*;
import com.example.ordermgmt.nexus.handler.ShippingService;
import io.temporal.api.enums.v1.ParentClosePolicy;
import io.temporal.common.SearchAttributeKey;
import io.temporal.common.converter.EncodedValues;
import io.temporal.failure.ActivityFailure;
import io.temporal.failure.ApplicationFailure;
import io.temporal.spring.boot.WorkflowImpl;
import io.temporal.workflow.*;
import org.slf4j.Logger;

import java.time.Duration;
import java.util.ArrayList;
import java.util.List;

@WorkflowImpl(taskQueues = "${orders.task-queue}")
public class OrderWorkflowScenarios implements DynamicWorkflow {
    private static final String BUG = "OrderWorkflowRecoverableFailure";
    private static final String CHILD = "OrderWorkflowChildWorkflow";
    private static final String SIGNAL = "OrderWorkflowHumanInLoopSignal";
    private static final String UPDATE = "OrderWorkflowHumanInLoopUpdate";
    private static final String VISIBILITY = "OrderWorkflowAdvancedVisibility";
    private static final String NEXUS = "OrderWorkflowNexusOperation";

    private static final Logger log = Workflow.getLogger(OrderWorkflowScenarios.class);

    private static final SearchAttributeKey<String> ORDER_STATUS_SA =
            SearchAttributeKey.forKeyword("OrderStatus");
    private final OrderActivities activities = Workflow.newActivityStub(OrderActivities.class,
            OrderActivities.defaultActivityOptions);
    private final OrderActivities localActivities = Workflow.newLocalActivityStub(OrderActivities.class,
            OrderActivities.defaultLocalActivityOptions);

    ShippingService shippingService = Workflow.newNexusServiceStub(
            ShippingService.class,
            NexusServiceOptions.newBuilder()
                    .setOperationOptions(
                            NexusOperationOptions.newBuilder()
                                    .setScheduleToCloseTimeout(Duration.ofSeconds(30))
                                    .build())
                    .setEndpoint(OrderApplication.defaultShippingNexusEndpoint)
                    .build());

    private int progress = 0;
    private String updatedAddress = null;

    @Override
    public Object execute(EncodedValues args) {
        Workflow.registerListener(new OrderWorkflowDynamicListenerImpl());
        OrderInput input = args.get(0, OrderInput.class);
        String type = Workflow.getInfo().getWorkflowType();
        log.info("Dynamic Order workflow started, type = {}, orderId = {}", type, input.getOrderId());

        // Create a saga to manage order compensations
        Saga saga = new Saga(new Saga.Options.Builder().setParallelCompensation(false).build());

        // Get items
        List<OrderItem> orderItems = localActivities.getItems();

        updateProgress("Check Fraud", 0, 0);

        // Check fraud
        activities.checkFraud(input);

        updateProgress("Prepare Shipment", 25, 1);

        // Prepare shipment
        saga.addCompensation(activities::undoPrepareShipment, input);
        activities.prepareShipment(input);

        updateProgress("Charge Customer", 50, 1);

        // Charge customer
        try {
            saga.addCompensation(activities::undoChargeCustomer, input);
            activities.chargeCustomer(input, type);
        } catch (ActivityFailure af) {
            log.error("Failed to charge customer", af);
            saga.compensate();
            throw af;
        }

        updateProgress("Ship Order", 75, 3);

        if (BUG.equals(type)) {
            // Simulate bug
            log.info("Throwing Exception to simulate a bug.  FIX to resolve.");
            throw new RuntimeException("Simulated bug - fix me!");
        }

        if (SIGNAL.equals(type) || UPDATE.equals(type)) {
            // Await message to update address
            waitForUpdatedAddressOrTimeout(input);
        }

        // Ship order items
        List<Promise<String>> promiseList = new ArrayList<>();
        for (OrderItem orderItem : orderItems) {
            log.info("Shipping item: {}", orderItem.getDescription());
            promiseList.add(shipItemAsync(input, orderItem, type));
        }

        // Wait for all items to ship
        Promise.allOf(promiseList).get();

        updateProgress("Order Completed", 100, 0);

        // Generate trackingId
        String trackingId = Workflow.randomUUID().toString();
        return new OrderOutput(trackingId, input.getAddress());
    }

    private Promise<String> shipItemAsync(OrderInput input, OrderItem orderItem, String type) {
        Promise<String> promise;
        if (CHILD.equals(type)) {
            // execute an async child wf to ship the item
            ShippingWorkflow orderShippingChild = Workflow.newChildWorkflowStub(ShippingWorkflow.class,
                    ChildWorkflowOptions.newBuilder()
                            .setWorkflowId("shipment-" + input.getOrderId() + "-" + orderItem.getId())
                            .setParentClosePolicy(ParentClosePolicy.PARENT_CLOSE_POLICY_TERMINATE)
                            .build());
            promise = Async.function(orderShippingChild::execute, new ShippingInput(input, orderItem));
        } else if (NEXUS.equals(type)) {
            // execute an async nexus operation to ship the item
            NexusOperationHandle<String> handle = Workflow.startNexusOperation(
                    shippingService::execute,
                    new ShippingInput(input, orderItem)
            );
            // wait for operation to start, then get promise for the result
            handle.getExecution().get();
            promise = handle.getResult();
        } else {
            // execute an async activity to ship the item
            promise = Async.function(activities::shipOrder, input, orderItem);
        }
        return promise;
    }

    private void waitForUpdatedAddressOrTimeout(OrderInput input) {
        log.info("Waiting up to 60 seconds for updated address");
        boolean ok = Workflow.await(Duration.ofSeconds(60), () -> updatedAddress != null);
        if (ok) {
            input.setAddress(updatedAddress);
        } else {
            // Do nothing - use the original address
            // In other cases, you may want to throw an exception on timeout, e.g.
            //   throw ApplicationFailure.newFailure (
            //     "Updated address was not received within 60 seconds.", "timeout");
        }
    }

    private void updateProgress(String orderStatus, int progress, int sleep) {
        this.progress = progress;
        if (sleep > 0) {
            Workflow.sleep(Duration.ofSeconds(sleep));
        }
        if (VISIBILITY.equals(Workflow.getInfo().getWorkflowType())) {
            Workflow.upsertTypedSearchAttributes(ORDER_STATUS_SA.valueSet(orderStatus));
        }
    }

    class OrderWorkflowDynamicListenerImpl implements OrderWorkflowMessages {

        @Override
        public int queryProgress() {
            return progress;
        }

        @Override
        public void updateOrderSignal(UpdateOrderInput updateInput) {
            log.info("Received update order signal with address: {}", updateInput.getAddress());
            updatedAddress = updateInput.getAddress();
        }

        @Override
        public String updateOrderUpdate(UpdateOrderInput updateInput) {
            log.info("Received update order update with address: {}", updateInput.getAddress());
            updatedAddress = updateInput.getAddress();
            return "Updated address: " + updatedAddress;
        }

        @Override
        public void updateOrderValidator(UpdateOrderInput updateInput) {
            if (!Character.isDigit(updateInput.getAddress().charAt(0))) {
                log.info("Rejecting order update, invalid address: {}", updateInput.getAddress());
                throw ApplicationFailure.newFailure("Address must start with a digit", "invalid-address");
            }
            log.info("Order update address is valid: {}", updateInput.getAddress());
        }
    }
}
