package com.example.ordermgmt.nexus.handler;

import com.example.ordermgmt.model.ShippingInput;
import com.example.ordermgmt.workflows.ShippingWorkflow;
import io.nexusrpc.handler.OperationHandler;
import io.nexusrpc.handler.OperationImpl;
import io.nexusrpc.handler.ServiceImpl;
import io.temporal.client.WorkflowOptions;
import io.temporal.nexus.Nexus;
import io.temporal.nexus.WorkflowRunOperation;
import io.temporal.spring.boot.NexusServiceImpl;
import org.springframework.stereotype.Component;


@Component
@NexusServiceImpl(taskQueues = "${shipping.task-queue}")
@ServiceImpl(service = ShippingService.class)
public class ShippingServiceImpl {
    @OperationImpl
    public OperationHandler<ShippingInput, String> execute() {
        return WorkflowRunOperation.fromWorkflowMethod(
                (ctx, details, input) ->
                        Nexus.getOperationContext()
                                .getWorkflowClient()
                                .newWorkflowStub(
                                        ShippingWorkflow.class,
                                        WorkflowOptions.newBuilder()
                                                .setWorkflowId(
                                                        String.format("shipment-%s-%s",
                                                                input.getOrderInput().getOrderId(),
                                                                input.getOrderItem().getId()))
                                                .build())
                                ::execute);
    }
}
