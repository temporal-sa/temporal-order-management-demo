require 'spec_helper'
require_relative '../activities/charge_customer_activity'
require_relative '../activities/prepare_shipment_activity'
require_relative '../activities/check_fraud_activity'
require_relative '../activities/ship_order_activity'
require_relative '../activities/undo_prepare_shipment_activity'
require_relative '../activities/undo_charge_customer_activity'
require_relative '../activities/get_items_activity'
require_relative '../models/order_input'
require_relative '../models/order_item'
require 'temporalio/testing'

describe Activities::ChargeCustomerActivity do
  let(:activity) { described_class.new }
  let(:input) { double('input') }

  before do
    info = double('info', attempt: 1)
    context = double('context', info: info)
    allow(Temporalio::Activity::Context).to receive(:current).and_return(context)
  end

  it 'returns SUCCESS when no error' do
    expect(activity.execute(input, 'NoError')).to eq('SUCCESS')
  end

  it 'raises StandardError for API unavailable' do
    expect { activity.execute(input, Activities::ChargeCustomerActivity::ERROR_CHARGE_API_UNAVAILABLE) }.to raise_error(StandardError)
  end

  it 'raises ArgumentError for invalid card' do
    expect { activity.execute(input, Activities::ChargeCustomerActivity::ERROR_INVALID_CREDIT_CARD) }.to raise_error(Temporalio::Error::ApplicationError)
  end
end

describe Activities::PrepareShipmentActivity do
  let(:activity) { described_class.new }
  let(:input) { double('input') }

  it 'returns SUCCESS' do
    expect(activity.execute(input)).to eq('SUCCESS')
  end
end

describe Activities::CheckFraudActivity do
  let(:activity) { described_class.new }
  let(:input) { double('input') }

  it 'returns SUCCESS' do
    expect(activity.execute(input)).to eq('SUCCESS')
  end
end

describe Activities::ShipOrderActivity do
  let(:activity) { described_class.new }
  let(:input) do
    Models::ShippingInput.new(Models::OrderInput.new('1', 'addr'), Models::OrderItem.new(1, 'desc', 1), 1)
  end

  it 'returns SUCCESS and heartbeats progress each second' do
    heartbeats = []
    env = Temporalio::Testing::ActivityEnvironment.new(on_heartbeat: ->(details) { heartbeats << details.first })
    expect(env.run(activity, input)).to eq('SUCCESS')
    expect(heartbeats).to eq([1, 2])
  end
end

describe Activities::UndoPrepareShipmentActivity do
  let(:activity) { described_class.new }
  let(:input) { double('input') }

  it 'returns SUCCESS' do
    expect(activity.execute(input)).to eq('SUCCESS')
  end
end

describe Activities::UndoChargeCustomerActivity do
  let(:activity) { described_class.new }
  let(:input) { double('input') }

  it 'returns SUCCESS' do
    expect(activity.execute(input)).to eq('SUCCESS')
  end
end

describe Activities::GetItemsActivity do
  let(:activity) { described_class.new }

  it 'returns an array of items' do
    result = activity.execute
    expect(result).to be_an(Array)
    expect(result.first).to respond_to(:id)
  end
end 