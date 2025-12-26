require 'temporalio/client'
require 'temporalio/worker'
require 'logger'
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

class Worker
  def run
    logger.info('Ruby order management worker starting...')
    Temporalio::Worker.new(
      client: temporal_client,
      task_queue: task_queue,
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
      ],
      workflow_payload_codec_thread_pool: Temporalio::Worker::ThreadPool.default
    ).run
  end

  private

  def logger
    @logger ||= Logger.new($stdout, level: :info)
  end

  def temporal_client
    @client ||= begin
      options = {
        logger: logger
      }.tap do |options|
        options.merge!(tls_options) if using_tls?
      end
      logger.info("Connecting to Temporal at #{temporal_address}")
      logger.info("Using namespace #{temporal_namespace}")
      logger.info("Using task queue #{task_queue}")
      logger.info("Using api key: #{api_key}")

      Temporalio::Client.connect(temporal_address, temporal_namespace, **options)
    end
  end

  def tls_options
    if api_key
      {

        tls: true,
        api_key: api_key,
        rpc_metadata: { 'temporal-namespace' => temporal_namespace }
      }
    elsif cert_path && key_path
      {
        tls: Temporalio::Client::Connection::TLSOptions.new(
          client_cert: File.read(cert_path),
          client_private_key: File.read(key_path)
        )
      }
    else
      {}
    end
  end

  def temporal_address
    ENV.fetch('TEMPORAL_ADDRESS', 'localhost:7233')
  end

  def temporal_namespace
    ENV.fetch('TEMPORAL_NAMESPACE', 'default')
  end

  def task_queue
    ENV.fetch('TEMPORAL_TASK_QUEUE', 'orders')
  end

  def api_key
    ENV['TEMPORAL_API_KEY']
  end

  def cert_path
    ENV['TEMPORAL_TLS_CLIENT_CERT_PATH']
  end

  def key_path
    ENV['TEMPORAL_TLS_CLIENT_KEY_PATH']
  end

  def using_tls?
    api_key || (cert_path && key_path)
  end
end

Worker.new.run if __FILE__ == $PROGRAM_NAME
