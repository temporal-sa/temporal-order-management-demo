require 'json/add/struct'
require_relative 'serialization'

module Models
  OrderOutput = Struct.new(:tracking_id, :address) do
    include Serialization
  end
end
