require 'json/add/struct'
require_relative 'serialization'

module Models
  UpdateOrderInput = Struct.new(:address) do
    include Serialization
  end
end
