"use client"

import { AddOfferFormData, AddOfferFormState } from "@/app/lib/definitions";
import React, { useActionState } from "react";
import { getAvailableModels } from "../../(home)/producers-and-models";
import { addOffer } from "@/app/actions/add-offer";

export default function AddOfferForm({ inputsData } : { inputsData : AddOfferFormData}) {
    const [availableModels, setAvailableModels] = React.useState<string[]>([]);

    function handleProducerChange(producer: string) {
        const models = getAvailableModels({ fieldName: "Producers", options: inputsData.producers, selected: [producer] }, inputsData.models);
        setAvailableModels(models);
    }

    const [isAuction, setIsAuction] = React.useState<boolean>(false);

    function changeOfferType(isAuction: boolean) {
        setIsAuction(isAuction);
    }

    const initialState: AddOfferFormState = {
        errors: {},
        values: {}
    };

    const [state, action] = useActionState(addOffer, initialState);

    const handleSubmit = (formData: FormData) => {
        formData.append("isAuction", isAuction.toString());
        const registraction_day = formData.get("dateOfFirstRegistration-day");
        const registraction_month = formData.get("dateOfFirstRegistration-month");
        const registraction_year = formData.get("dateOfFirstRegistration-year");
        if (registraction_day && registraction_month && registraction_year) {
            const date = `${registraction_year.toString().padStart(4, "0")}-${registraction_month.toString().padStart(2, "0")}-${registraction_day.toString().padStart(2, "0")}`;
            formData.set("dateOfFirstRegistration", date);
            formData.delete("dateOfFirstRegistration-day");
            formData.delete("dateOfFirstRegistration-month");
            formData.delete("dateOfFirstRegistration-year");
        }
        const auction_day = formData.get("auctionEndDate-day");
        const auction_month = formData.get("auctionEndDate-month");
        const auction_year = formData.get("auctionEndDate-year");
        if (auction_day && auction_month && auction_year) {
            const auctionDate = `${auction_year.toString().padStart(4, "0")}-${auction_month.toString().padStart(2, "0")}-${auction_day.toString().padStart(2, "0")}`;
            formData.set("auctionEndDate", auctionDate);
            formData.delete("auctionEndDate-day");
            formData.delete("auctionEndDate-month");
            formData.delete("auctionEndDate-year");
        }
        return action(formData)
    }

    // UI

    function SelectionLabel({ id, name, options }: { id: string, name: string; options: string[] }) {
        return (
            <>
                <label htmlFor={id} className="text-lg font-semibold">
                    {name}
                </label>
                <select id={id} name={id} className="border rounded p-2" required defaultValue="">
                    <option value="" disabled>Select a {name}...</option>
                    {options.map((item: string) => (
                        <option key={item} value={item}>{item}</option>
                    ))}
                </select>
            </>
        )
    }

    function NumberInputField({ id, name } : { id:string, name: string }) {
        return (
            <>
                <label htmlFor={id} className="text-lg font-semibold">
                    {name}
                </label>
                <input
                    type="number"
                    id={id}
                    name={id}
                    className="border rounded p-2"
                    placeholder={`Enter ${name.toLowerCase()}`}
                    required
                />
            </>
        )
    }

    function DateSelectionField({ id, name } : { id: string, name: string }) {
        return (
            <>
                <label htmlFor={id} className="text-lg font-semibold">
                    {name}
                </label>
                <div className="flex gap-2 items-end">

                    <div className="flex flex-col">
                        <input
                            type="number"
                            id={`${id}-day`}
                            name={`${id}-day`}
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
                            id={`${id}-month`}
                            name={`${id}-month`}
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
                            id={`${id}-year`}
                            name={`${id}-year`}
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

    function TextInputField({ id, name } : { id: string, name: string }) {
        return (
            <>
                <label htmlFor={id} className="text-lg font-semibold">
                    {name}
                </label>
                <input
                    type="text"
                    id={id}
                    name={id}
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
                <SelectionLabel id="model" name="Model" options={availableModels} />
                <SelectionLabel id="color" name="Color" options={inputsData.colors} />
                <SelectionLabel id="gearbox" name="Gearbox" options={inputsData.gearboxes} />
                <SelectionLabel id="fuelType" name="Fuel Type" options={inputsData.fuelTypes} />
                <SelectionLabel id="driveType" name="Drive Type" options={inputsData.driveTypes} />
                <SelectionLabel id="country" name="Country" options={inputsData.countries} />
                <NumberInputField id="productionYear" name="Production Year" />
                <NumberInputField id="mileage" name="Mileage" />
                <NumberInputField id="numberOfDoors" name="Number of doors" />
                <NumberInputField id="numberOfSeats" name="Number of seats" />
                <NumberInputField id="power" name="Power" />
                <NumberInputField id="engineDisplacement" name="Engine displacement" />
                <DateSelectionField id="dateOfFirstRegistration" name="Date of first registration" />
                <TextInputField id="plateNumber" name="Plate number" />
                <TextInputField id="location" name="Location" />

                <label htmlFor="description" className="text-lg font-semibold">
                    Description
                </label>
                <textarea
                    id="description"
                    name="description"
                    className="border rounded p-2 h-32"
                    required
                ></textarea>

                <label htmlFor="images" className="text-lg font-semibold">
                    Images
                </label>
                //TODO
                {/* <input type="file" id="images" name="images" className="border rounded p-2" multiple required /> */}

                <div className="my-10"/>
                <NumberInputField id="price" name="Price" />
                <label htmlFor="offer type" className="text-lg font-semibold">
                    Offer Type
                </label>
                <div className="flex gap-8">
                    <div className="flex items-center">
                        <input
                            id="standard offer"
                            type="radio"
                            value="P"
                            checked={isAuction === false}
                            onChange={() => changeOfferType(false)}
                            className="h-4 w-4 text-blue-500 focus:ring-blue-400"
                        />
                        <label htmlFor="standard offer" className="ml-2 text-sm font-medium text-gray-900">
                            Standard offer
                        </label>
                    </div>
                    <div className="flex items-center">
                        <input
                            id="auction"
                            type="radio"
                            value="C"
                            checked={isAuction === true}
                            onChange={() => changeOfferType(true)}
                            className="h-4 w-4 text-blue-500 focus:ring-blue-400"
                        />
                        <label htmlFor="auction" className="ml-2 text-sm font-medium text-gray-900">
                            Auction
                        </label>
                    </div>
                </div>

                {isAuction && (
                    <>
                        <DateSelectionField id="auctionEndDate" name="Auction end date" />
                        <NumberInputField id="buyNowAuctionPrice" name="Buy now price" />
                    </>
                )}

                <button type="submit" className="bg-blue-600 text-white rounded p-2">Add offer</button>
            </div>
        </form>
    )
}