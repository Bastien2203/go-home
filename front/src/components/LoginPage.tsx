import { useEffect, useState, type ChangeEvent, type FormEvent } from "react"
import { api } from "../services/api";


export const LoginPage = (props: {
    onLoginSuccess: () => void
}) => {
    const [canRegister, setCanRegister] = useState<boolean>()
    const [message, setMessage] = useState<string>()
    const [formData, setFormData] = useState<{
        email?: string;
        password?: string
    }>({})

    useEffect(() => {
        api.canRegister()
        .then(r => setCanRegister(r.can_register))
        .catch(_ => setMessage("network error"))
    }, [props.onLoginSuccess])

    const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        if (!formData.email || !formData.password) {
            setMessage("Please fill all fields")
            return
        }

        const subscriptionMethod = canRegister ? api.register : api.login
        subscriptionMethod(formData.email, formData.password)
            .then(props.onLoginSuccess)
            .catch(_ => setMessage("error invalid login"))
    }

    const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
        setFormData(prev => ({
            ...prev,
            [e.target.id]: e.target.value
        }))
    }

    return <div className="flex min-h-[80vh] w-full items-center justify-center bg-gray-50 px-4">
    <div className="w-full max-w-sm rounded-xl border border-gray-100 bg-white p-8 shadow-lg">
        
        <div className="mb-6 text-center">
            <h2 className="text-2xl font-bold tracking-tight text-gray-900">
                {canRegister ? "Create an account" : "Welcome back"}
            </h2>
            <p className="mt-2 text-sm text-gray-500">
                {canRegister ? "Enter your details below" : "Please enter your details"}
            </p>
        </div>

        <form id="loginForm" onSubmit={handleSubmit} className="space-y-5">
            <div>
                <label htmlFor="email" className="mb-2 block text-sm font-medium text-gray-700">
                    Email address
                </label>
                <input 
                    type="email" 
                    id="email" 
                    placeholder="name@company.com" 
                    required 
                    onChange={handleChange} 
                    className="w-full rounded-lg border border-gray-300 px-4 py-2.5 text-gray-900 placeholder:text-gray-400 focus:border-blue-600 focus:outline-none focus:ring-1 focus:ring-blue-600 sm:text-sm"
                />
            </div>

            <div>
                <label htmlFor="password" className="mb-2 block text-sm font-medium text-gray-700">
                    Password
                </label>
                <input 
                    type="password" 
                    id="password" 
                    placeholder="••••••••" 
                    required 
                    onChange={handleChange} 
                    className="w-full rounded-lg border border-gray-300 px-4 py-2.5 text-gray-900 placeholder:text-gray-400 focus:border-blue-600 focus:outline-none focus:ring-1 focus:ring-blue-600 sm:text-sm"
                />
            </div>

            <button 
                type="submit"
                className="w-full rounded-lg bg-blue-600 px-5 py-2.5 text-center text-sm font-medium text-white transition-colors hover:bg-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-300"
            >
                {canRegister ? "Sign Up" : "Sign In"}
            </button>
        </form>

        {message && (
            <div className={`mt-4 rounded-md p-3 text-sm ${message.includes('error') ? 'bg-red-50 text-red-700' : 'bg-green-50 text-green-700'}`}>
                <p className="text-center font-medium">{message}</p>
            </div>
        )}
    </div>
</div>
}