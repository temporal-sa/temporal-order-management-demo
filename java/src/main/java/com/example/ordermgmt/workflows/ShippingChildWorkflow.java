package com.example.ordermgmt.workflows;


import com.example.ordermgmt.model.OrderInput;
import com.example.ordermgmt.model.OrderItem;
import io.temporal.workflow.WorkflowInterface;
import io.temporal.workflow.WorkflowMethod;


@WorkflowInterface
public interface ShippingChildWorkflow {
    @WorkflowMethod
    void execute(OrderInput input, OrderItem item);
}
