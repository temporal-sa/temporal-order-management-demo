from flask import Flask, render_template, request, jsonify
import uuid
import os
from client import get_client
from temporalio.client import WorkflowFailureError
from data import OrderInput
import time

app = Flask(__name__)

# Sample data for the order
order_data = {
    "table_top": {"item": "Table Top", "quantity": 1},
    "table_legs": {"item": "Table Legs", "quantity": 2},
    "keypad": {"item": "Keypad", "quantity": 1},
}

# Sample data for payment
payment_data = {
    "name": "Billy Bob",
    "card_type": "Visa",
    "card_number": "1234567890",
}

# Sample data for shipping
shipping_data = {
    "name": "Billy Bob",
    "address": "12345 Dongle Way, Nowhere California",
}

# Sample choices for the drop-down menu
scenarios = [
    "HappyPath",
    "AdvancedVisibility",
    "HumanInLoopSignal",
    "HumanInLoopUpdate",
    "ChildWorkflow",
    "APIFailure",
    "RecoverableFailure",
    "NonRecoverableFailure",
]

@app.route('/', methods=['GET', 'POST'])
async def main_order_page():
    order_id = str(uuid.uuid4().int)[:6] 
    return render_template('index.html', order_data=order_data, payment_data=payment_data, shipping_data=shipping_data, scenarios=scenarios, order_id=order_id)

@app.route('/process_order')
async def process_order():
    selected_scenario = request.args.get('scenario')
    order_id = request.args.get('order_id')
    client = await get_client()

    input = OrderInput(
        OrderId= order_id,
        Address=shipping_data["address"],
    )               

    await client.start_workflow(
        "OrderWorkflow"+selected_scenario,
        input,
        id=f'order-{order_id}',
        task_queue=os.getenv("TEMPORAL_TASK_QUEUE"),
    )    

    return render_template('process_order.html', selected_scenario=selected_scenario, oder_id=order_id)

@app.route('/order_confirmation')
async def order_confirmation():
    order_id = request.args.get('order_id')

    client = await get_client()
    order_workflow = client.get_workflow_handle(f'order-{order_id}')
    order_output = await order_workflow.result()
    
    tracking_id = order_output["trackingId"]
    address = order_output["address"]

    return render_template('order_confirmation.html', order_id=order_id, tracking_id=tracking_id, address=address)

@app.route('/get_progress')
async def get_progress():
    order_id = request.args.get('order_id')

    try:
        client = await get_client()
        order_workflow = client.get_workflow_handle(f'order-{order_id}')
        progress_percent = await order_workflow.query("getProgress")

        return jsonify({"progress": progress_percent})
    except Exception as e:
        print(f"Error in get_progress route: {str(e)}")
        return jsonify({"error": "Internal Server Error"}), 500

if __name__ == '__main__':
    app.run(debug=True)    
