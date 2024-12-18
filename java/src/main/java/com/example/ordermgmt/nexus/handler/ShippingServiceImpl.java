package com.example.ordermgmt.nexus.handler;

import com.example.ordermgmt.model.OrderInput;
import com.example.ordermgmt.model.OrderItem;
import com.example.ordermgmt.model.OrderOutput;
import com.example.ordermgmt.model.ShippingInput;
import com.example.ordermgmt.nexus.service.ShippingService;
import com.example.ordermgmt.nexus.workflows.ShippingWorkflow;
import io.nexusrpc.handler.OperationHandler;
import io.nexusrpc.handler.OperationImpl;
import io.nexusrpc.handler.ServiceImpl;
import io.temporal.client.WorkflowOptions;
import io.temporal.nexus.WorkflowClientOperationHandlers;

@ServiceImpl(service = ShippingService.class)
public class ShippingServiceImpl {
    @OperationImpl
    public OperationHandler<ShippingInput, OrderOutput> execute () {
        return WorkflowClientOperationHandlers.fromWorkflowMethod(
                (ctx, details, client, input) ->
                        client.newWorkflowStub(
                                ShippingWorkflow.class,
                                // Workflow IDs should typically be business meaningful IDs and are used to
                                // dedupe workflow starts.
                                // For this example, we're using the request ID allocated by Temporal when the
                                // caller workflow schedules
                                // the operation, this ID is guaranteed to be stable across retries of this
                                // operation.
                                //
                                // Task queue defaults to the task queue this operation is handled on.
                                WorkflowOptions.newBuilder().setWorkflowId(details.getRequestId()).build())
                                ::execute);
    }
}
