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

@WorkflowImpl(taskQueues = "${orders.task-queue}")
class OrderWorkflowImpl implements OrderWorkflow {
    private static final Logger log = Workflow.getLogger(OrderWorkflowImpl.class);

    private final OrderActivities activities = Workflow.newActivityStub(
            OrderActivities.class,
            OrderActivities.defaultActivityOptions
    );

    private final OrderActivities localActivities = Workflow.newLocalActivityStub(
            OrderActivities.class,
            OrderActivities.defaultLocalActivityOptions
    );

    private int progress = 0;

    @Override
    public OrderOutput execute(OrderInput input) {
        String type = Workflow.getInfo().getWorkflowType();
        log.info("Order workflow started, type = {}, orderId = {}", type, input.getOrderId());

        // Get items
        List<OrderItem> orderItems = localActivities.getItems();

        // Check fraud
        activities.checkFraud(input);
        sleep(1, 25);

        // Prepare shipment
        activities.prepareShipment(input);
        sleep(1, 50);

        // Charge customer
        activities.chargeCustomer(input, type);
        sleep(3, 75);

        // Ship order items
        List<Promise<String>> promiseList = new ArrayList<>();
        for (OrderItem orderItem : orderItems) {
            log.info("Shipping item: {}", orderItem.getDescription());
            promiseList.add(Async.function(activities::shipOrder, input, orderItem));
        }

        // Wait for all items to ship
        Promise.allOf(promiseList).get();
        sleep(0, 100);

        // Generate trackingId
        String trackingId = Workflow.randomUUID().toString();
        return new OrderOutput(trackingId, input.getAddress());
    }

    @Override
    public int queryProgress() {
        return progress;
    }

    private void sleep(int sleep, int progress) {
        this.progress = progress;
        if (sleep > 0) {
            Workflow.newTimer(
                    Duration.ofSeconds(sleep),
                    TimerOptions.newBuilder()
                            .setSummary("Sleep")
                            .build()
            ).get();
        }
    }
}
