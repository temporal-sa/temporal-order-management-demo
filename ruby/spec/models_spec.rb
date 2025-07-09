require 'spec_helper'
require_relative '../models/order_input'
require_relative '../models/order_output'
require_relative '../models/order_item'
require_relative '../models/update_order_input'
require_relative '../models/activity_response'

describe Models::OrderInput do
  it 'initializes with order_id and address' do
    input = described_class.new(1, 'addr')
    expect(input.order_id).to eq(1)
    expect(input.address).to eq('addr')
  end
end

describe Models::OrderOutput do
  it 'initializes with tracking_id and address' do
    output = described_class.new('track', 'addr')
    expect(output.tracking_id).to eq('track')
    expect(output.address).to eq('addr')
  end
end

describe Models::OrderItem do
  it 'initializes with id, description, and quantity' do
    item = described_class.new(1, 'desc', 2)
    expect(item.id).to eq(1)
    expect(item.description).to eq('desc')
    expect(item.quantity).to eq(2)
  end
end

describe Models::UpdateOrderInput do
  it 'initializes with address' do
    input = described_class.new('addr')
    expect(input.address).to eq('addr')
  end
end

describe Models::ActivityResponse do
  it 'initializes with order_id and status' do
    resp = described_class.new(1, 'ok')
    expect(resp.order_id).to eq(1)
    expect(resp.status).to eq('ok')
  end
end

describe 'Serialization' do
  it 'deep_camelize_keys works for OrderInput' do
    input = Models::OrderInput.new(1, 'addr')
    camel = input.deep_camelize_keys
    expect(camel).to include('orderId', 'address')
  end
end 