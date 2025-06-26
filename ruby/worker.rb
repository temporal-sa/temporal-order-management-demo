require 'temporalio/client'
require_relative 'shared_objects'
require_relative 'activities/order_activities'
require_relative 'workflows/order_workflow'
require_relative 'workflows/order_workflow_scenarios'
require_relative 'workflows/shipping_child_workflow'

if __FILE__ == $0
  address = ENV['TEMPORAL_ADDRESS'] || '127.0.0.1:7233'
  namespace = ENV['TEMPORAL_NAMESPACE'] || 'default'
  task_queue = ENV['TEMPORAL_TASK_QUEUE'] || 'orders'
  client = Temporalio::Client.connect(address, namespace)
  activities = OrderActivities.new
  # The following is a placeholder; actual worker registration may differ depending on SDK maturity
  # You may need to use Temporalio::Worker or similar if available in the gem
  puts "Connected to Temporal on #{address}"
  puts "Ruby order management worker ready (manual workflow/activity registration may be required)"
  # TODO: Register workflows and activities with the worker as per SDK docs when available
end
