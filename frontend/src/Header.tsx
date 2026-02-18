import { FaBars } from "react-icons/fa6";
function Header({roomName}: {roomName:string}) {
	return (
		<div className="navbar bg-base-100 shadow-sm">
			<div className="flex-none">
				<button className="btn btn-square btn-ghost">
				  <FaBars />
				</button>

			</div>
		  <a className="btn btn-ghost text-xl gap-10 align-center" href="/">🎉 vnc.party</a>
		  
			  <div className="text-2xl align-right bold font-mono">
			    {roomName}
			  </div>
		</div>
	)
}

export default Header;
