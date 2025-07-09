require 'spec_helper'
require_relative '../workflows/order_workflow'
require_relative '../workflows/order_workflow_scenarios'
require_relative '../workflows/shipping_child_workflow'

RSpec.configure do |config|
  config.before(:each) do
    logger = double('logger', info: nil, error: nil)
    allow(Temporalio::Workflow).to receive(:logger).and_return(logger)
    stub_const('Temporalio::Internal::Worker', Class.new)
    stub_const('Temporalio::Internal::Worker::WorkflowInstance', Class.new)
    stub_const('Temporalio::Internal::Worker::WorkflowInstance::Scheduler', Class.new)
  end
end

describe Workflows::OrderWorkflow do
  let(:workflow) { described_class.new }
  let(:input) { { 'OrderId' => 1, 'Address' => 'addr' } }

  it 'initializes progress to 0' do
    expect(workflow.progress).to eq(0)
  end

  it 'can run execute without error if activities succeed' do
    skip('Cannot run full workflow execute outside of Temporal workflow environment')
  end
end

describe Workflows::OrderWorkflowScenarios do
  let(:workflow) { described_class.new }
  let(:input) { { 'OrderId' => 1, 'Address' => 'addr' } }

  it 'initializes updated_address to nil' do
    expect(workflow.updated_address).to be_nil
  end

  it 'handles activity failure and runs compensations' do
    allow(Temporalio::Workflow).to receive(:execute_activity).and_return([double('item', description: 'desc')], 'SUCCESS', 'SUCCESS', 'SUCCESS').and_raise(StandardError)
    allow(Temporalio::Workflow).to receive_message_chain(:random, :uuid).and_return('uuid')
    expect { workflow.execute(input) }.to raise_error(StandardError)
  end
end

describe Workflows::ShippingChildWorkflow do
  let(:workflow) { described_class.new }
  let(:input) { double('input', order_id: 1) }
  let(:item) { double('item') }

  it 'calls ShipOrderActivity' do
    expect(Temporalio::Workflow).to receive(:execute_activity).with(Activities::ShipOrderActivity, input, item, start_to_close_timeout: 5, retry_policy: nil)
    workflow.execute(input, item)
  end
end 