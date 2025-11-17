import { useState } from "react";

import { api } from "../services/api";
import type { Device, Parser } from "../types";
import { useParsers } from "../hooks/useParsers";

const CreateDeviceModal = (props: {
    refreshDeviceList: () => void,
    deviceTypes: Record<string, string> | undefined
}) => {
    const [form, setForm] = useState<Device>({
        name: "",
        addr: "",
        type: "temperature_sensor",
        parser_type: "",
    });
    const { parsers } = useParsers()

    if (!parsers) return <div>Loading parsers...</div>;
    if (!props.deviceTypes) return <div>Loading device types...</div>;


    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        setForm({ ...form, [e.target.name]: e.target.value });
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        api.createDevice(form).then(props.refreshDeviceList)
    };

    return (
        <div className="p-4 bg-white rounded-xl shadow-md max-w-md mx-auto">
            <h2 className="text-lg font-semibold mb-4">Create Device</h2>

            <form onSubmit={handleSubmit} className="space-y-3">
                <input
                    name="name"
                    placeholder="Device name"
                    value={form.name}
                    onChange={handleChange}
                    className="border p-2 w-full rounded"
                />
                <input
                    name="addr"
                    placeholder="Device address"
                    value={form.addr}
                    onChange={handleChange}
                    className="border p-2 w-full rounded"
                />

                <select
                    name="parser_type"
                    value={form.parser_type}
                    onChange={handleChange}
                    className="border p-2 w-full rounded"
                >
                    <option value="">Select parser</option>
                    {parsers?.map((p: Parser) => (
                        <option key={p.name} value={p.name}>
                            {p.name}
                        </option>
                    ))}
                </select>

                <select
                    name="type"
                    value={form.type}
                    onChange={handleChange}
                    className="border p-2 w-full rounded"
                >
                    {
                        Object.entries(props.deviceTypes).map(e => {
                            const type = e[0]
                            const name = e[1]
                            return <option value={type} key={type}>{name}</option>
                        })
                    }

                </select>

                <button
                    type="submit"

                    className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 w-full cursor-pointer"
                >
                    Create
                </button>
            </form>
        </div>
    );
};

export default CreateDeviceModal;
