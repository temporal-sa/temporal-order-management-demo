package com.example.shippingmgmt;

import com.example.ordermgmt.nexus.handler.ShippingServiceImpl;
import com.example.ordermgmt.nexus.workflows.ShippingWorkflowImpl;
import com.example.ordermgmt.workflows.OrderWorkflowScenarios;
import com.example.ordermgmt.workflows.ShippingChildWorkflow;
import com.example.ordermgmt.workflows.ShippingChildWorkflowImpl;
import io.grpc.netty.shaded.io.netty.handler.ssl.SslContext;
import io.grpc.netty.shaded.io.netty.handler.ssl.SslContextBuilder;
import io.grpc.util.AdvancedTlsX509KeyManager;
import io.temporal.api.cloud.cloudservice.v1.DeleteServiceAccountRequest;
import io.temporal.client.WorkflowClient;
import io.temporal.client.WorkflowClientOptions;
import io.temporal.serviceclient.SimpleSslContextBuilder;
import io.temporal.serviceclient.WorkflowServiceStubs;
import io.temporal.serviceclient.WorkflowServiceStubsOptions;
import io.temporal.worker.*;
import io.temporal.workflow.NexusServiceOptions;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.ApplicationContext;
import org.springframework.context.annotation.Bean;
import org.springframework.core.env.Environment;
import org.springframework.stereotype.Component;
import io.grpc.netty.shaded.io.grpc.netty.GrpcSslContexts;

import java.io.File;
import java.io.FileInputStream;
import java.lang.management.ManagementFactory;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.Collections;

@Slf4j
@SpringBootApplication
public class ShippingApplication {


    public static void main(String[] args) {
        SpringApplication.run(ShippingApplication.class, args);
    }

    @Component
    public static class StartupRunner implements CommandLineRunner {
        @Value("${ordermgmt.task-queue}")
        String taskQueue;
        @Value("${ordermgmt.task-queue-shipping}")
        String shippingTaskQueue;
        @Value("${nexus.worker.client-cert-file:unspecified}")
        String clientCertFileValue;
        @Value("${nexus.worker.client-key-file:unspecified}")
        String clientKeyFileValue;
        @Value("${nexus.worker.target-endpoint:unspecified}")
        String nexusTargetEndpointValue;
        @Value("${nexus.worker.namespace:unspecified}")
        String nexusNamespaceValue;
        @Value("${nexus.worker.refresh-period:30}")
        String refreshPeriodValue;
        @Value("${nexus.worker.endpoint:unspecified}")
        String shippingEndpointValue;

        @Autowired
        private WorkerFactory workerFactory;

        //@Autowired
        //private ApplicationContext context;  // Added to allow for activity registration of additional worker for nexus purposes

        @Autowired
        private Environment env;

        @Override
        public void run(String... args) throws Exception {
            log.info("Getting worker for task queue '{}'", taskQueue);
            Worker worker = workerFactory.getWorker(taskQueue);


            log.info("Registering dynamic workflow on a worker with task queue '{}'", taskQueue);
            worker.registerWorkflowImplementationTypes(ShippingChildWorkflowImpl.class);  // Register child so it does not run the dynamic workflow
            worker.registerWorkflowImplementationTypes(OrderWorkflowScenarios.class);

            log.info("*** Endpoint is set to  [{}]", shippingEndpointValue);
            worker.registerWorkflowImplementationTypes(
                    WorkflowImplementationOptions.newBuilder()
                            .setNexusServiceOptions(
                                    Collections.singletonMap(
                                            "ShippingService",
                                            NexusServiceOptions.newBuilder().setEndpoint(shippingEndpointValue).build()
                                    )
                            )
                            .build(),
                    ShippingWorkflowImpl.class
            );
           // startNexusWorker();

            log.info("Starting worker factory");
            workerFactory.start();
        }


    private void startNexusWorker() throws Exception {
        WorkflowServiceStubs service;
        log.info("The environment being used is [{}]", env.getActiveProfiles());
        if ((env.getActiveProfiles().length > 0) && (env.getActiveProfiles()[0].equals("tc"))) {
            log.info("Using tc profile  - setting up the ssl connectivity");

            File clientCertFile = new File(clientCertFileValue);
            File clientKeyFile = new File(clientKeyFileValue);
            long refreshPeriod = refreshPeriodValue != null ? Integer.parseInt(refreshPeriodValue) : 0;

            SslContext sslContext = SimpleSslContextBuilder.forPKCS8(
                            new FileInputStream(clientCertFile), new FileInputStream(clientKeyFile))
                    .build();

            if (refreshPeriod > 0) {
                AdvancedTlsX509KeyManager clientKeyManager = new AdvancedTlsX509KeyManager();
                // Reload credentials every minute
                clientKeyManager.updateIdentityCredentialsFromFile(
                        clientKeyFile,
                        clientCertFile,
                        refreshPeriod,
                        TimeUnit.MINUTES,
                        Executors.newScheduledThreadPool(1));
                sslContext =
                        GrpcSslContexts.configure(SslContextBuilder.forClient().keyManager(clientKeyManager))
                                .build();
            }

            log.info("THe nexus target endpoint for the Nexus worker is [{}]", nexusTargetEndpointValue);
            service = WorkflowServiceStubs.newServiceStubs(
                            WorkflowServiceStubsOptions.newBuilder()
                                    .setSslContext(sslContext)
                                    .setTarget(nexusTargetEndpointValue)
                                    .build());


        } else {
            // The service is without any sslContext (Running locally without auth.)
            service = WorkflowServiceStubs.newServiceStubs(
                            WorkflowServiceStubsOptions.newBuilder()
                                    .setTarget(nexusTargetEndpointValue)
                                    .build());
        }
        WorkflowClient client =
                WorkflowClient.newInstance(
                        service, WorkflowClientOptions.newBuilder().setNamespace(nexusNamespaceValue).build());

        // worker factory that can be used to create workers for specific task queues
        WorkerFactory factory = WorkerFactory.newInstance(client);
        WorkerOptions workerOptions = WorkerOptions.newBuilder()
                                          .setIdentity(shippingEndpointValue + "-" + ManagementFactory.getRuntimeMXBean().getName())
                                          .build();
        Worker shippingWorker = factory.newWorker(shippingTaskQueue, workerOptions);
        shippingWorker.registerActivitiesImplementations(context.getBean("orderActivitiesImpl"));
        shippingWorker.registerWorkflowImplementationTypes(ShippingWorkflowImpl.class);
        shippingWorker.registerNexusServiceImplementation(new ShippingServiceImpl());

        factory.start();
    }
}
/**
    @Bean
    public CommandLineRunner commandLineRunner(ApplicationContext ctx) {
        return args -> {

            System.out.println("Let's inspect the beans provided by Spring Boot:");

            String[] beanNames = ctx.getBeanDefinitionNames();
            //Arrays.sort(beanNames);
            for (String beanName : beanNames) {
                System.out.println(beanName);
            }
        };
    }
    **/
}
