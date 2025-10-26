export default function LoginPage() {
    return (
		<div className={`flex flex-col items-center justify-center min-h-screen bg-zinc-900 text-white`}>
			<div className={`w-full max-w-md p-8 mx-auto bg-zinc-800 rounded-lg shadow-lg`}>
				<h1 className={`text-3xl font-bold text-center mb-6`}>Login</h1>
				<div className={`h-1 my-3`}></div>
				<p className={"text-center text-zinc-400 mb-8"}>
                    Access your account
                </p>
                <form>
                    <div className="mb-4">
                        <label className="block text-zinc-400 mb-2" htmlFor="email">
                            Email
                        </label>
                        <input
                            id="email"
                            type="email"
                            className="w-full p-3 bg-zinc-700 rounded border border-zinc-600 focus:outline-none focus:ring-2 focus:ring-pink-500"
                            required
                            placeholder="you@example.com"
                        />
                    </div>
                    <div className="mb-4">
                        <label className="block text-zinc-400 mb-2" htmlFor="password">
                            Password
                        </label>
                        <input
                            id="password"
                            type="password"
                            className="w-full p-3 bg-zinc-700 rounded border border-zinc-600 focus:outline-none focus:ring-2 focus:ring-pink-500"
                            required
                            placeholder="*******"
                        />
                    </div>
                    <div>
                        <button
                            type="submit"
                            className="w-full bg-pink-600 hover:bg-pink-700 text-white font-bold py-3 px-4 rounded-lg focus:outline-none focus:shadow-outline transition duration-300"
                        >
                            Login
                        </button>
                    </div>
                </form>
			</div>
		</div>
	)
}