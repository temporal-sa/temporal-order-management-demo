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

    public static void main(String[] args) {
        SpringApplication.run(OrderApplication.class, args);
    }

    @Component
    public static class StartupRunner implements CommandLineRunner {
        @Value("${ordermgmt.task-queue}")
        String taskQueue;

        @Autowired
        private WorkerFactory workerFactory;

        @Override
        public void run(String... args) throws Exception {
            log.info("Getting worker for task queue '{}'", taskQueue);
            Worker worker = workerFactory.getWorker(taskQueue);

            log.info("Registering dynamic workflow on a worker with task queue '{}'", taskQueue);
            worker.registerWorkflowImplementationTypes(OrderWorkflowScenarios.class);

            log.info("Starting worker factory");
            workerFactory.start();
        }
    }
}
