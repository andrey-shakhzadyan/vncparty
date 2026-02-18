import { useForm, Form } from "react-hook-form"
import { useNavigate } from "react-router";
//import { ErrorMessage } from "@hookform/error-message"
//import { useState } from "react";

interface VNCForm {
	server_addr: string
	server_port: number
	room_name: string
}

function LandingContent() {
	const {
		register,
		control,
		formState: { isSubmitSuccessful, errors }
	} = useForm<VNCForm>();

	let navigate = useNavigate();
	//  const [redirectId, setRedirectId] = useState("");
	// const onSubmit: SubmitHandler<VNCForm> = (data) => { // change to form component

	// 	fetch("/create_room", {
	// 		method: "POST",
	// 		headers: {
	// 			'Accept': 'application/json',
	// 			'Content-Type': 'application/json'
	// 		},
	// 		body: JSON.stringify(data)
	// 	})
	// 		.then((response) => response.json())
	// 		.then((respjson) => console.log(respjson))
	// 		.catch((error) => console.log("Error: " + error));
	// }; 	<form className="fieldset" onSubmit={handleSubmit(onSubmit)}>
	if (isSubmitSuccessful) {
		console.log("submitted")
	}
	return (
		<div className="hero min-h-screen">
			<div className="hero-content flex-row text-center content-center">
				<h1 className="font-bold text-9xl"> 🎉 vnc.party </h1>
				<div className="card w-full max-w-sm shrink-0 shadow-2xl">
					<div className="card-body">
						<Form
							action="/api/create_room"
							control={control}
							onSuccess={async (resp) => {
								var rjson = await resp.response.json()
								navigate("/room?uuid=" + rjson.room_id)

							}}
							onError={() => {
								console.log("form submission failed")
							}}
							headers={{
								"Accept": "application/json",
								"Content-Type": "application/json"
							}}
						>
							<div className="join">
								<label className="label">Server address</label>
							</div>
							<div className="join">
								<input type="text" className="input join-item" placeholder="Server address" {...register("server_addr", { required: true })} ></input>
								<input type="number" className="input join-item w-1/5 no-arrow" placeholder="15901" {...register("server_port", { valueAsNumber: true, required: false })} ></input>
							</div>
							{errors.server_addr && <div className="text-red-600">Server IP is required</div>}
							<label className="label">Room Name</label>
						  <input type="text" className="input" placeholder="Room Name" {...register("room_name", { required: true, maxLength: 20, pattern: /^[A-Za-z0-9]+$/i })} ></input>
						  {errors.room_name && <div className="text-red-600">Room name is required, and must be an alphanumeric string of length less than 20</div>}

							<div className="divider"></div>
							<button className="btn btn-primary">Create room & share</button>
						</Form>
					</div>
				</div>
			</div>
		</div >
	)
}

export default LandingContent;
