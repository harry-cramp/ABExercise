import { useState, useEffect } from 'react'
import './App.css'

type InventoryResponse = {
	quantity: number;
};

type BuyResponse = {
	ticket: number;
};

type BuyErrorResponse = {
	message: string;
};

const ItemPage: React.FC = () => {
	const [loading, setLoading] = useState<boolean | null>(true);
	const [quantity, setQuantity] = useState<number | null>(null);
	const [ticketNumber, setTicketNumber] = useState<number | null>(null);
	const [errorMsg, setErrorMsg] = useState<string | null>(null);
	
	useEffect(() => {
		const fetchQuantity = async() => {
			try {
				const response = await fetch("http://localhost:8080/inventory");
				
				if(!response.ok) {
					throw new Error("Failed to retrieve item quantity");
				}
				
				const data: InventoryResponse = await response.json();
				setLoading(false);
				setQuantity(data.quantity);
				
				if(data.quantity === 0) {
					setErrorMsg("Out of Stock");
				}
			}catch(err) {
				setErrorMsg("Error: failed to retrieve inventory");
			}
		};
		
		fetchQuantity();
	}, []);
	
	const buyItem = async() => {
		if(loading || quantity == 0) {
			return;
		}
	
		setLoading(true);
		
		try {
			const idemKey = crypto.randomUUID();
			const response = await fetch("http://localhost:8080/buy", {
				method: "POST",
				headers: {
					"Idempotency-Key": idemKey,
				}
			});
			
			if(!response.ok) {
				const errorData: BuyErrorResponse = await response.json();

		    throw new Error(errorData.message);
			}
			
			const data: BuyResponse = await response.json();
			setLoading(false);
			setTicketNumber(data.ticket);
		}catch(err) {
			setErrorMsg((err as Error).message);
		}
	};
	
	return (
		<div className="itemPage">
      <div className="itemImage">
      </div>
      <div className="itemInfo">
        <div className="itemTitle">
          New Phone
        </div>
        {ticketNumber == null && <div className="itemQuantity">
          Quantity: {loading ? "Retrieving quantity..." : quantity}
        </div>}
        {errorMsg === null && ticketNumber !== null && <div className="itemTicket">
        	<span style={{"color": "green"}}>Success!</span> Thank you, your ticket number is: {ticketNumber}
        </div>}
        {errorMsg !== null && <div style={{"color": "red"}}>
        	{errorMsg}
        </div>}
	      <button className={loading || quantity == 0 ? "itemBuyButtonDisabled" : ""} onClick={() => { buyItem(); }} disabled={loading}>
	        Buy Now
	      </button>
      </div>
    </div>
	);
};

function App() {
  return (
    <>
      <title>
        We Sell New Phones
      </title>
      <ItemPage/>
    </>
  );
}

export default App
