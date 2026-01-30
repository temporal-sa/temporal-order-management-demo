require_relative 'activities/get_items_activity'
require_relative 'activities/check_fraud_activity'
require_relative 'activities/prepare_shipment_activity'
require_relative 'activities/charge_customer_activity'
require_relative 'activities/ship_order_activity'
require_relative 'activities/undo_prepare_shipment_activity'
require_relative 'activities/undo_charge_customer_activity'
require_relative 'workflows/order_workflow'
require_relative 'workflows/order_workflow_scenarios'
require_relative 'workflows/shipping_child_workflow'
require 'logger'
require 'temporalio/client'
require 'temporalio/env_config'
require 'temporalio/worker'

args, kwargs = Temporalio::EnvConfig::ClientConfig.load_client_connect_options
client = Temporalio::Client.connect(
  *args, **kwargs,
  logger: Logger.new($stdout, level: Logger::INFO)
)
puts "✅ Client connected to #{args[0]} in namespace '#{args[1]}'"

task_queue = ENV.fetch('TEMPORAL_TASK_QUEUE', 'orders')
worker = Temporalio::Worker.new(
  client:,
  task_queue:,
  workflows: [
    Workflows::OrderWorkflow,
    Workflows::OrderWorkflowScenarios,
    Workflows::ShippingChildWorkflow
  ],
  activities: [
    Activities::GetItemsActivity,
    Activities::CheckFraudActivity,
    Activities::PrepareShipmentActivity,
    Activities::ChargeCustomerActivity,
    Activities::ShipOrderActivity,
    Activities::UndoPrepareShipmentActivity,
    Activities::UndoChargeCustomerActivity
  ]
)

# Run the worker until SIGINT
puts 'Starting worker (ctrl+c to exit)'
worker.run(shutdown_signals: ['SIGINT'])
