package com.example.ordermgmt;


import com.example.ordermgmt.workflows.ShippingWorkflow;
import com.example.ordermgmt.workflows.OrderWorkflowScenarios;

import io.temporal.worker.*;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.core.env.Environment;
import org.springframework.stereotype.Component;


@Slf4j
@SpringBootApplication
public class OrderApplication {


    public static void main(String[] args) {
        SpringApplication.run(OrderApplication.class, args);
    }

    @Component
    public static class StartupRunner implements CommandLineRunner {
        @Value("${ordermgmt.task-queue}")
        String taskQueue;

        @Autowired
        private WorkerFactory workerFactory;

        @Autowired
        private Environment env;

        @Override
        public void run(String... args) throws Exception {
            log.info("Getting worker for task queue '{}'", taskQueue);
            Worker worker = workerFactory.getWorker(taskQueue);

            log.info("Registering dynamic workflow on a worker with task queue '{}'", taskQueue);
            worker.registerWorkflowImplementationTypes(ShippingWorkflow.ShippingWorkflowImpl.class);
            worker.registerWorkflowImplementationTypes(OrderWorkflowScenarios.class);

            log.info("Starting worker factory");
            workerFactory.start();
        }
}
}
