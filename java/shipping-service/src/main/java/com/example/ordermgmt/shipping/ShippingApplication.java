package com.example.ordermgmt.shipping;

import com.example.ordermgmt.nexus.handler.ShippingServiceImpl;
import io.temporal.worker.Worker;
import io.temporal.worker.WorkerFactory;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.stereotype.Component;

@Slf4j
@ComponentScan({"com.example.ordermgmt.activities"})
@SpringBootApplication
public class ShippingApplication {

    public static void main(String[] args) {
        SpringApplication.run(ShippingApplication.class, args);
    }

    @Component
    public static class StartupRunner implements CommandLineRunner {
        @Value("${shipping.task-queue}")
        String workerTaskQueue;

        @Autowired
        private WorkerFactory workerFactory;

        @Override
        public void run(String... args) throws Exception {
            log.info("Getting worker for task queue '{}'", workerTaskQueue);
            Worker worker = workerFactory.getWorker(workerTaskQueue);

            // auto discovery via @NexusServiceImpl did not work, so manually registering here. need to revisit
            log.info("Registering nexus service implementation on a worker with task queue '{}'", workerTaskQueue);
            worker.registerNexusServiceImplementation(new ShippingServiceImpl());

            log.info("Starting worker factory");
            workerFactory.start();
        }
    }
}
