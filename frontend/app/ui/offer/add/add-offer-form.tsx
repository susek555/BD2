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

    // UI

    function SelectionLabel({ name, options }: { name: string; options: string[] }) {
        return (
            <>
                <label htmlFor={`${name.toLowerCase()}`} className="text-lg font-semibold">
                    {name}
                </label>
                <select id={`${name.toLowerCase()}`} name={`${name.toLowerCase()}`} className="border rounded p-2" required defaultValue="">
                    <option value="" disabled>Select a {name}...</option>
                    {options.map((item: string) => (
                        <option key={item} value={item}>{item}</option>
                    ))}
                </select>
            </>
        )
    }

    function NumberInputField({ name } : { name: string }) {
        return (
            <>
                <label htmlFor={`${name.toLowerCase()}`} className="text-lg font-semibold">
                    {name}
                </label>
                <input
                    type="number"
                    id={`${name.toLowerCase()}`}
                    name={`${name.toLowerCase()}`}
                    className="border rounded p-2"
                    placeholder={`Enter ${name.toLowerCase()}`}
                    required
                />
            </>
        )
    }

    function DateSelectionField({ name } : { name: string }) {
        return (
            <>
                <label htmlFor={`${name.toLowerCase()}`} className="text-lg font-semibold">
                    {name}
                </label>
                <div className="flex gap-2 items-end">

                    <div className="flex flex-col">
                        <input
                            type="number"
                            id={`${name.toLowerCase()}-day`}
                            name={`${name.toLowerCase()}-day`}
                            min={1}
                            max={31}
                            className="border rounded p-2 w-20"
                            required
                            placeholder="Day"
                        />
                    </div>
                    <div className="flex flex-col">
                        <input
                            type="number"
                            id={`${name.toLowerCase()}-month`}
                            name={`${name.toLowerCase()}-month`}
                            min={1}
                            max={12}
                            className="border rounded p-2 w-20"
                            required
                            placeholder="Month"
                        />
                    </div>
                    <div className="flex flex-col">
                        <input
                            type="number"
                            id={`${name.toLowerCase()}-year`}
                            name={`${name.toLowerCase()}-year`}
                            min={1900}
                            max={2100}
                            className="border rounded p-2 w-28"
                            required
                            placeholder="Year"
                        />
                    </div>
                </div>
            </>
        )
    }

    function TextInputField({ name } : { name: string }) {
        return (
            <>
                <label htmlFor={`${name.toLowerCase()}`} className="text-lg font-semibold">
                    {name}
                </label>
                <input
                    type="text"
                    id={`${name.toLowerCase()}`}
                    name={`${name.toLowerCase()}`}
                    className="border rounded p-2"
                    placeholder={`Enter ${name.toLowerCase()}`}
                    required
                />
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
                <NumberInputField name="Production Year" />
                <NumberInputField name="Mileage" />
                <NumberInputField name="Number of doors" />
                <NumberInputField name="Number of seats" />
                <NumberInputField name="Power" />
                <NumberInputField name="Engine displacement" />
                <DateSelectionField name="Date of first registration" />
                <TextInputField name="Plate number" />
                <TextInputField name="Location" />

                <label htmlFor="description" className="text-lg font-semibold">Description</label>
                <textarea id="description" name="description" className="border rounded p-2 h-32" required></textarea>

                <label htmlFor="images" className="text-lg font-semibold">Images</label>
                //TODO
                {/* <input type="file" id="images" name="images" className="border rounded p-2" multiple required /> */}

                <div className="my-10"/>
                <NumberInputField name="Price" />



                <button type="submit" className="bg-blue-600 text-white rounded p-2">Submit</button>
            </div>
        </form>
    )
}