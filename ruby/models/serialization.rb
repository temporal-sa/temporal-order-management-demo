require 'json'

module Models
  module Serialization
    def self.included(base)
      base.class_eval do
        def deep_camelize_keys
          deep_to_h.deep_camelize_keys
        end

        def deep_to_h
          hash = {}
          members.each do |member|
            value = send(member)
            hash[member] = if value.respond_to?(:deep_to_h)
                             value.deep_to_h
                           elsif value.respond_to?(:to_h)
                             value.to_h
                           else
                             value
                           end
          end
          hash
        end
      end
    end
  end
end

class Hash
  def deep_camelize_keys
    transform_keys { |key| key.to_s.gsub(/_([a-z])/) { $1.upcase } }.transform_values do |value|
      case value
      when Hash
        value.deep_camelize_keys
      when Array
        value.map { |item| item.is_a?(Hash) ? item.deep_camelize_keys : item }
      else
        value
      end
    end
  end
end
