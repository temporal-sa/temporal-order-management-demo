package com.example.ordermgmt;

import com.example.ordermgmt.workflows.OrderWorkflowScenarios;
import io.temporal.worker.Worker;
import io.temporal.worker.WorkerFactory;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.stereotype.Component;

@Slf4j
@SpringBootApplication
public class OrderApplication {
    public static final String defaultShippingNexusEndpoint =
            System.getenv("TEMPORAL_NEXUS_SHIPPING_ENDPOINT") != null ?
                    System.getenv("TEMPORAL_NEXUS_SHIPPING_ENDPOINT") :
                    "shipping-endpoint";

    public static void main(String[] args) {
        SpringApplication.run(OrderApplication.class, args);
    }

    @Component
    public static class StartupRunner implements CommandLineRunner {
        @Value("${orders.task-queue}")
        String workerTaskQueue;

        @Autowired
        private WorkerFactory workerFactory;

        @Override
        public void run(String... args) throws Exception {
            log.info("Getting worker for task queue '{}'", workerTaskQueue);
            Worker worker = workerFactory.getWorker(workerTaskQueue);

            log.info("Registering dynamic workflow on a worker with task queue '{}'", workerTaskQueue);
            worker.registerWorkflowImplementationTypes(OrderWorkflowScenarios.class);

            // unfortunately, the below does not work, as WorkflowImplementationOptions are not applied to
            // DynamicWorkflows. (this would work if OrderWorkflowScenarios was a regular workflow)
//            worker.registerWorkflowImplementationTypes(
//                    WorkflowImplementationOptions.newBuilder()
//                            .setNexusServiceOptions(
//                                    Collections.singletonMap(
//                                            "ShippingService",
//                                            NexusServiceOptions.newBuilder()
//                                                    .setEndpoint(defaultShippingNexusEndpoint)
//                                                    .build()))
//                            .build(),
//                    OrderWorkflowScenarios.class);

            log.info("Starting worker factory");
            workerFactory.start();
        }
    }
}
