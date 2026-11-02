using DotNetOrderManagement.model;
using NexusRpc;

namespace DotNetOrderManagement;

// Nexus service contract for shipping. The names must match the handler
// (see go/app/service.go): service "shipping-service", operation "ship-item".
[NexusService("shipping-service")]
public interface IShippingService
{
    [NexusOperation("ship-item")]
    string ShipItem(ShippingInput input);
}
