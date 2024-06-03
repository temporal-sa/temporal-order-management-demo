package com.example.ordermgmt.workflows;

import com.example.ordermgmt.activities.OrderActivities;
import com.example.ordermgmt.model.OrderInput;
import com.example.ordermgmt.model.OrderItem;
import io.temporal.spring.boot.WorkflowImpl;
import io.temporal.workflow.Workflow;
import io.temporal.workflow.WorkflowInterface;
import io.temporal.workflow.WorkflowMethod;
import org.slf4j.Logger;

@WorkflowInterface
public interface ShippingChildWorkflow {
    @WorkflowMethod
    void execute(OrderInput input, OrderItem item);
}

@WorkflowImpl(taskQueues = "${ordermgmt.task-queue}")
class ShippingChildWorkflowImpl implements ShippingChildWorkflow {
    private static final Logger log = Workflow.getLogger(ShippingChildWorkflowImpl.class);

    private final OrderActivities activities = Workflow.newActivityStub(OrderActivities.class,
            OrderActivities.defaultActivityOptions);

    @Override
    public void execute(OrderInput input, OrderItem item) {
        log.info("Shipping workflow started, orderId = {}", input.getOrderId());

        // Ship order
        activities.shipOrder(input, item);
    }
}
