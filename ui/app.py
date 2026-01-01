from quart import Quart, render_template, request, jsonify, Response
import uuid
import os
import asyncio
import json
from client import get_client
from data import OrderInput, UpdateOrder

app = Quart(__name__)

client = None

@app.before_serving
async def startup():
    global client
    client = await get_client()

# Order Info
order_data = {
    "table_top": {"item": "Table Top", "quantity": 1},
    "table_legs": {"item": "Table Legs", "quantity": 2},
    "keypad": {"item": "Keypad", "quantity": 1},
}

# Payment Info
payment_data = {
    "name": "Alice Jones",
    "card_type": "Visa",
    "card_number": "1111222233334444",
}

# Shipping Info
shipping_data = {
    "name": "Alice Jones",
    "address": "123 Main St. Redwood, CA",
}

# Scenario choices dropdown
scenarios = [
    "HappyPath",
    "AdvancedVisibility",
    "HumanInLoopSignal",
    "HumanInLoopUpdate",
    "ChildWorkflow",
    "NexusOperation",
    "APIFailure",
    "RecoverableFailure",
    "NonRecoverableFailure",
]

api_key = os.getenv("TEMPORAL_API_KEY")

if api_key:
    scenarios.append("APIKeyRotation")

@app.route('/', methods=['GET', 'POST'])
async def main_order_page():
    order_id = str(uuid.uuid4().int)[:6]
    return await render_template(
        'index.html',
        order_data=order_data,
        payment_data=payment_data,
        shipping_data=shipping_data,
        scenarios=scenarios,
        order_id=order_id,
        api_key=api_key,
    )

@app.route('/process_order')
async def process_order():
    selected_scenario = request.args.get('scenario')
    order_id = request.args.get('order_id')

    input = OrderInput(
        OrderId= order_id,
        Address=shipping_data["address"],
    )

    await client.start_workflow(
        "OrderWorkflow"+selected_scenario,
        input,
        id=f'order-{order_id}',
        task_queue=os.getenv("TEMPORAL_TASK_QUEUE", "orders"),
    )

    return await render_template('process_order.html', selected_scenario=selected_scenario, order_id=order_id)

@app.route('/order_confirmation')
async def order_confirmation():
    order_id = request.args.get('order_id')

    order_workflow = client.get_workflow_handle(f'order-{order_id}')
    order_output = await order_workflow.result()

    tracking_id = order_output["trackingId"]
    address = order_output["address"]

    return await render_template('order_confirmation.html', order_id=order_id, tracking_id=tracking_id, address=address)

@app.route('/stream_progress')
async def stream_progress():
    order_id = request.args.get('order_id')

    async def event_stream():
        while True:
            try:
                order_workflow = client.get_workflow_handle(f'order-{order_id}')
                progress_percent = await order_workflow.query("getProgress")

                desc = await order_workflow.describe()
                if desc.status == 3:
                    error_message = f"Workflow failed: order-{order_id}"
                    print(f"Error in stream_progress route: {error_message}")
                    yield f"data: {json.dumps({'error': error_message})}\n\n"
                    break

                yield f"data: {json.dumps({'progress': progress_percent})}\n\n"

                if progress_percent >= 100:
                    break

                await asyncio.sleep(1)
            except Exception as e:
                print(f"Error in stream_progress: {str(e)}")
                yield f"data: {json.dumps({'error': str(e)})}\n\n"
                break

    headers = {
        "Cache-Control": "no-cache",
        "Connection": "keep-alive",
		"X-Accel-Buffering": "no",
    }
    return Response(event_stream(), headers=headers, content_type="text/event-stream")

@app.route('/signal', methods=['POST'])
async def signal():
    order_id = request.args.get('order_id')
    data = await request.get_json()
    address = data.get('address')

    SignalOrderInput = UpdateOrder(
        Address=address
    )

    try:
        order_workflow = client.get_workflow_handle(f'order-{order_id}')
        await order_workflow.signal("UpdateOrder", SignalOrderInput)
    except Exception as e:
        print(f"Error sending signal: {str(e)}")
        return jsonify({"error": str(e)}), 500

    return 'Signal received successfully', 200

@app.route('/update', methods=['POST'])
async def update():
    order_id = request.args.get('order_id')
    data = await request.get_json()
    address = data.get('address')

    UpdateOrderInput = UpdateOrder(
        Address=address
    )

    update_result = None
    try:
        order_workflow = client.get_workflow_handle(f'order-{order_id}')
        update_result = await order_workflow.execute_update(
            update="UpdateOrder",
            arg=UpdateOrderInput,
        )
    except Exception as e:
        print(f"Error sending signal: {str(e)}")
        result = f"Update for order_id {order_id} rejected, not a valid address! {str(e)}"
        return jsonify(result=result)

    result = f"Update for order_id {order_id} accepted: {update_result}"

    return jsonify(result=result)

if __name__ == '__main__':
    app.run(debug=True)
