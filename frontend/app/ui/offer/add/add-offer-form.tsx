"use client"

import { AddOfferFormData } from "@/app/lib/definitions";
import React from "react";
import { getAvailableModels } from "../../(home)/producers-and-models";

export default function AddOfferForm({ inputsData } : { inputsData : AddOfferFormData}) {

    function handleSubmit(formData: FormData): any {
        const formDataObj = Object.fromEntries(formData.entries());
        console.log("Add Offer form data:", formDataObj);
        return formDataObj;
    }

    const [availableModels, setAvailableModels] = React.useState<string[]>([]);

    function handleProducerChange(producer: string) {
        const models = getAvailableModels({ fieldName: "Producers", options: inputsData.producers, selected: [producer] }, inputsData.models);
        setAvailableModels(models);
    }

    function SelectionLabel({ name, options }: { name: string; options: string[] }) {
        return (
            <>
                <label htmlFor={`${name.toLowerCase()}`} className="text-lg font-semibold">{name}</label>
                <select id={`${name.toLowerCase()}`} name={`${name.toLowerCase()}`} className="border rounded p-2" required defaultValue="">
                    <option value="" disabled>Select a {name}...</option>
                    {options.map((item: string) => (
                        <option key={item} value={item}>{item}</option>
                    ))}
                </select>
            </>
        )
    }

    return (
        <form className=" w-full md:w-200" action={handleSubmit}>
            <div className="rounded-lg bg-gray-50 px-6 pb-4 pt-8 flex flex-col gap-4">
                <label htmlFor="producer" className="text-lg font-semibold">Producer</label>
                <select
                    id="producer"
                    name="producer"
                    className="border rounded p-2"
                    required
                    defaultValue=""
                    onChange={e => handleProducerChange(e.target.value)}
                >
                    <option value="" disabled>Select a producer...</option>
                    {inputsData.producers.map((item: string) => (
                        <option key={item} value={item}>{item}</option>
                    ))}
                </select>
                <SelectionLabel name="Model" options={availableModels} />
                <SelectionLabel name="Color" options={inputsData.colors} />
                <SelectionLabel name="Gearbox" options={inputsData.gearboxes} />
                <SelectionLabel name="Fuel Type" options={inputsData.fuelTypes} />
                <SelectionLabel name="Drive Type" options={inputsData.driveTypes} />
                <SelectionLabel name="Country" options={inputsData.countries} />

                <label htmlFor="description" className="text-lg font-semibold">Description</label>
                <textarea id="description" name="description" className="border rounded p-2 h-32" required></textarea>

                <button type="submit" className="bg-blue-600 text-white rounded p-2">Submit</button>
            </div>
        </form>
    )
}