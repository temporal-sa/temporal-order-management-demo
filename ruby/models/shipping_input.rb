require 'json/add/struct'
require_relative 'serialization'

module Models
  # index: item n ships in n*2 seconds; nil means a random delay
  ShippingInput = Struct.new(:order, :item, :index) do
    include Serialization
  end
end
