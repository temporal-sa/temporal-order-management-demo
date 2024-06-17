package com.example.ordermgmt.workflows;

import com.example.ordermgmt.activities.OrderActivities;
import com.example.ordermgmt.model.OrderInput;
import com.example.ordermgmt.model.OrderItem;
import com.example.ordermgmt.model.OrderOutput;
import io.temporal.spring.boot.WorkflowImpl;
import io.temporal.workflow.*;
import org.slf4j.Logger;

import java.time.Duration;
import java.util.ArrayList;
import java.util.List;

@WorkflowInterface
public interface OrderWorkflow {
    @WorkflowMethod(name = "OrderWorkflowHappyPath")
    OrderOutput execute(OrderInput input);

    @QueryMethod(name = "getProgress")
    int queryProgress();
}

@WorkflowImpl(taskQueues = "${ordermgmt.task-queue}")
class OrderWorkflowImpl implements OrderWorkflow {
    private static final Logger log = Workflow.getLogger(OrderWorkflowImpl.class);

    private final OrderActivities activities = Workflow.newActivityStub(OrderActivities.class,
            OrderActivities.defaultActivityOptions);

    private final OrderActivities localActivities = Workflow.newLocalActivityStub(OrderActivities.class,
            OrderActivities.defaultLocalActivityOptions);

    private int progress = 0;

    @Override
    public OrderOutput execute(OrderInput input) {
        String type = Workflow.getInfo().getWorkflowType();
        log.info("Order workflow started, type ={}, orderId = {}", type, input.getOrderId());

        // Get items
        List<OrderItem> orderItems = localActivities.getItems();

        // Check fraud
        activities.checkFraud(input);
        updateProgress(25, 1);

        // Prepare shipment
        activities.prepareShipment(input);
        updateProgress(50, 1);

        // Charge customer
        activities.chargeCustomer(input, type);
        updateProgress(75, 3);

        // Ship orders
        List<Promise<Void>> promiseList = new ArrayList<>();
        for (OrderItem orderItem : orderItems) {
            log.info("Shipping item: {}", orderItem.getDescription());
            promiseList.add(Async.procedure(activities::shipOrder, input, orderItem));
        }

        // Wait for all items to ship
        Promise.allOf(promiseList).get();
        updateProgress(100, 1);

        // Generate trackingId
        String trackingId = Workflow.randomUUID().toString();
        return new OrderOutput(trackingId, input.getAddress());
    }

    @Override
    public int queryProgress() {
        return progress;
    }

    private void updateProgress(int progress, int sleep) {
        this.progress = progress;
        if (sleep > 0) {
            Workflow.sleep(Duration.ofSeconds(sleep));
        }
    }
}

