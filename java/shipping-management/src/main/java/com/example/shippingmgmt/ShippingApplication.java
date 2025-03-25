package com.example.shippingmgmt;

import com.example.ordermgmt.customize.TemporalOptionsConfig;
import com.example.ordermgmt.nexus.handler.ShippingServiceImpl;
import com.example.ordermgmt.workflows.ShippingWorkflow;
import io.temporal.worker.*;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.ApplicationContext;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Import;
import org.springframework.core.env.Environment;
import org.springframework.stereotype.Component;

@Slf4j
@ComponentScan({"com.example.ordermgmt"})
@SpringBootApplication
@Import({TemporalOptionsConfig.class})
public class ShippingApplication {


    public static void main(String[] args) {
        SpringApplication.run(ShippingApplication.class, args);
    }

    @Component
    public static class StartupRunner implements CommandLineRunner {
        @Value("${ordermgmt.task-queue}")
        String taskQueue;

        @Autowired
        private WorkerFactory workerFactory;

        @Autowired
        private ApplicationContext context;  // Added to allow for activity registration of additional worker for nexus purposes

        @Autowired
        private Environment env;

        @Override
        public void run(String... args) throws Exception {
            log.info("Getting worker for task queue '{}'", taskQueue);
            Worker worker = workerFactory.getWorker(taskQueue);

            worker.registerWorkflowImplementationTypes(ShippingWorkflow.ShippingWorkflowImpl.class);
            worker.registerNexusServiceImplementation(new ShippingServiceImpl());
            log.info("Starting worker factory");
            workerFactory.start();
        }

}

}
