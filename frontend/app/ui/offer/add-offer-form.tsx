"use client"

import { OfferFormData, OfferFormState } from "@/app/lib/definitions/offer-form";
import React, { useActionState, useEffect } from "react";
import { getAvailableModels } from "@/app/ui/(filters-sidebar)/producers-and-models";
import { addOffer } from "@/app/actions/add-offer";
import { editOffer } from "@/app/actions/edit-offer";
import { parseOfferForm, OfferFormEnum } from "@/app/lib/utils";

export const offerActionEnum = {
    ADD_OFFER: true,
    EDIT_OFFER: false
}

export function OfferForm(
{
    inputsData,
    initialValues = {is_auction: false},
    apiAction = offerActionEnum.ADD_OFFER,
    id = undefined
} : {
    inputsData : OfferFormData
    initialValues?: Partial<OfferFormState['values']>,
    apiAction?: boolean
    id?: string | undefined
}) {
    // Validate that ID is provided when edit mode is active
    if (apiAction === offerActionEnum.EDIT_OFFER && id === undefined) {
        throw new Error("ID is required when editing an offer");
    }


    const initialState: OfferFormState = {
        errors: {},
        values: initialValues as OfferFormState['values'],
    };

    const offerWrapper = (state: OfferFormState, formData: FormData) => {
            const detailsPart = formData.get('detailsPart') === 'true';
            formData.delete('detailsPart');

            const { boolean: result, offerFormState } = parseOfferForm(formData, detailsPart);

            // If there are validation errors or pricing is to be set, return the new state without calling API
            if (result === OfferFormEnum.pricingPartLeft) {
                return offerFormState;
            }

            // Otherwise call the API action with validated data
            if (apiAction) {
                return addOffer(offerFormState);
            } else {
                return editOffer(offerFormState, id!);
            }
    };

    const [state, action] = useActionState(offerWrapper, initialState);

    useEffect(() => {
        if (state.values?.producer) {
            setProducer(state.values.producer.toString());
            handleProducerChange(state.values.producer.toString());
        }
        if (typeof state.values?.is_auction === "boolean") {
            setIsAuction(state.values.is_auction);
        }

        console.log("State errors:", state.errors);
        console.log("State values:", state.values);

        if (Object.keys(state?.errors || {}).length === 0 && state.values?.producer) {
            setDetailsPart(false);
        }
    }
    , [state]);

    const [detailsPart, setDetailsPart] = React.useState<boolean>(true);
    const [producer, setProducer] = React.useState<string>(state.values?.producer?.toString() ?? "");

    const [availableModels, setAvailableModels] = React.useState<string[]>([]);

    const handleProducerChange = (newProducer: string) => {
        const models = getAvailableModels(
            { fieldName: "Producers", options: inputsData.producers, selected: [newProducer] },
            inputsData.models
        );
        setAvailableModels(models);
    };

    const [is_auction, setIsAuction] = React.useState<boolean>(state.values?.is_auction ?? false);

    const changeOfferType = (isAuction: boolean) => {
        setIsAuction(isAuction);
    };

    const handleSubmit = (formData: FormData) => {

        // formData.append("isAuction", liveState.values?.isAuction == null ? "false" : liveState.values.isAuction.toString());
        const registraction_day = formData.get("registration_date-day");
        const registraction_month = formData.get("registration_date-month");
        const registraction_year = formData.get("registration_date-year");
        if (registraction_day && registraction_month && registraction_year) {
            const date = `${registraction_year.toString().padStart(4, "0")}-${registraction_month.toString().padStart(2, "0")}-${registraction_day.toString().padStart(2, "0")}`;
            formData.set("registration_date", date);
            formData.delete("registration_date-day");
            formData.delete("registration_date-month");
            formData.delete("registration_date-year");
            // updateField("registration_date", date);
        }
        const auction_day = formData.get("auction_end_date-day");
        const auction_month = formData.get("auction_end_date-month");
        const auction_year = formData.get("auction_end_date-year");
        if (auction_day && auction_month && auction_year) {
            const auctionDate = `${auction_year.toString().padStart(4, "0")}-${auction_month.toString().padStart(2, "0")}-${auction_day.toString().padStart(2, "0")}`;
            formData.set("auction_end_date", auctionDate);
            formData.delete("auction_end_date-day");
            formData.delete("auction_end_date-month");
        // Add detailsPart to formData
        }
        formData.append('detailsPart', detailsPart.toString());

        setProducer(formData.get("producer")?.toString() ?? "");
        const isAuctionValue = formData.get("is_auction") === "true";
        setIsAuction(isAuctionValue);
        formData.set("is_auction", isAuctionValue.toString());

        action(formData)
    }

    // UI

    function SelectionLabel({ id, name, options, required = true }: {
        id: string,
        name: string;
        options: string[]
        required?: boolean;
    }) {
        return (
            <>
                <label htmlFor={id} className="text-lg font-semibold">
                    {name}
                </label>
                <select
                    id={id}
                    name={id}
                    className="border rounded p-2"
                    required={required}
                    defaultValue={
                        typeof state?.values?.[id] === "string"
                            ? state.values[id]
                            : ""
                    }
                >
                    <option value="" disabled>Select a {name}...</option>
                    {options.map((item: string) => (
                        <option key={item} value={item}>{item}</option>
                    ))}
                </select>
                <div id="username-error" aria-live="polite" aria-atomic="true">
                    {state?.errors?.[id] &&
                        state.errors[id].map((error: string) => (
                        <p className="mt-2 text-sm text-red-500" key={error}>
                            {error}
                        </p>
                    ))}
                </div>
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
                    defaultValue={state?.values?.[id]?.toString() ?? "0"}
                    required
                />
                <div id="username-error" aria-live="polite" aria-atomic="true">
                    {state?.errors?.[id] &&
                        state.errors[id].map((error: string) => (
                        <p className="mt-2 text-sm text-red-500" key={error}>
                            {error}
                        </p>
                    ))}
                </div>
            </>
        )
    }

    function DateSelectionField({ id, name }: { id: string, name: string }) {
        const initialDate = state.values?.[id]?.toString() ?? "";
        const initialDay = initialDate ? initialDate.split("-")[2] ?? "" : "";
        const initialMonth = initialDate ? initialDate.split("-")[1] ?? "" : "";
        const initialYear = initialDate ? initialDate.split("-")[0] ?? "" : "";

        return (
            <>
                <label htmlFor={id} className="text-lg font-semibold">{name}</label>
                <div className="flex gap-2 items-end">
                    <input
                        type="number"
                        id={`${id}-day`}
                        name={`${id}-day`}
                        min={1}
                        max={31}
                        className="border rounded p-2 w-20"
                        defaultValue={initialDay}
                        required
                        placeholder="Day"
                    />
                    <input
                        type="number"
                        id={`${id}-month`}
                        name={`${id}-month`}
                        min={1}
                        max={12}
                        className="border rounded p-2 w-20"
                        defaultValue={initialMonth}
                        required
                        placeholder="Month"
                    />
                    <input
                        type="number"
                        id={`${id}-year`}
                        name={`${id}-year`}
                        min={1900}
                        max={2100}
                        className="border rounded p-2 w-28"
                        defaultValue={initialYear}
                        required
                        placeholder="Year"
                    />
                </div>
                <div id="username-error" aria-live="polite" aria-atomic="true">
                    {state?.errors?.[id]?.map((error: string) => (
                        <p className="mt-2 text-sm text-red-500" key={error}>{error}</p>
                    ))}
                </div>
            </>
        );
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
                    defaultValue={state?.values?.[id]?.toString() ?? ""}
                    placeholder={`Enter ${name.toLowerCase()}`}
                    required
                />
                <div id="username-error" aria-live="polite" aria-atomic="true">
                    {state?.errors?.[id] &&
                        state.errors[id].map((error: string) => (
                        <p className="mt-2 text-sm text-red-500" key={error}>
                            {error}
                        </p>
                    ))}
                </div>
            </>
        )
    }

    return (
        <form className=" w-full md:w-200" action={handleSubmit}>
            <div className={`rounded-lg bg-gray-50 px-6 pb-4 pt-8 flex flex-col gap-4 ${detailsPart ? "block" : "hidden"}`}>
                <>
                    <label htmlFor="producer" className="text-lg font-semibold">Producer</label>
                    <select
                        id="producer"
                        name="producer"
                        className="border rounded p-2"
                        required
                        // defaultValue={
                        //     typeof state?.values?.producer === "string"
                        //         ? state.values.producer
                        //         : ""
                        // }
                        value={producer}
                        onChange={e => {
                            handleProducerChange(e.target.value);
                            setProducer(e.target.value); // Aktualizuj lokalny stan
                        }}
                    >
                        <option value="" disabled>Select a producer...</option>
                        {inputsData.producers.map((item: string) => (
                            <option key={item} value={item}>{item}</option>
                        ))}
                    </select>
                    <div id="username-error" aria-live="polite" aria-atomic="true">
                        {state?.errors?.producer &&
                            state.errors.producer.map((error: string) => (
                            <p className="mt-2 text-sm text-red-500" key={error}>
                                {error}
                            </p>
                        ))}
                    </div>
                    {/* <SelectionLabel id="producer" name="Producer" options={inputsData.producers} /> */}
                    <SelectionLabel id="model" name="Model" options={availableModels} />
                    <SelectionLabel id="color" name="Color" options={inputsData.colors} />
                    <SelectionLabel id="fuel_type" name="Fuel Type" options={inputsData.fuelTypes} />
                    <SelectionLabel id="drive" name="Drive Type" options={inputsData.driveTypes} />
                    <SelectionLabel id="transmission" name="Transmission" options={inputsData.gearboxes} />
                    <NumberInputField id="number_of_gears" name="Number of gears" />
                    <NumberInputField id="production_year" name="Production Year" />
                    <NumberInputField id="mileage" name="Mileage" />
                    <NumberInputField id="number_of_doors" name="Number of doors" />
                    <NumberInputField id="number_of_seats" name="Number of seats" />
                    <NumberInputField id="engine_power" name="Power" />
                    <NumberInputField id="engine_capacity" name="Engine displacement" />
                    <DateSelectionField id="registration_date" name="Date of first registration" />
                    <TextInputField id="registration_number" name="Plate number" />
                    <TextInputField id="vin" name="VIN" />

                    <label htmlFor="description" className="text-lg font-semibold">
                        Description
                    </label>
                    <textarea
                        id="description"
                        name="description"
                        className="border rounded p-2 h-32"
                        defaultValue={state?.values?.description?.toString() ?? ""}
                        placeholder="Enter description..."
                        required
                    ></textarea>

                    <label htmlFor="images" className="text-lg font-semibold">
                        Images
                    </label>
                    TODO
                    {/* <input type="file" id="images" name="images" className="border rounded p-2" multiple required /> */}

                    <button type="submit" className="bg-blue-600 text-white rounded p-2">Set Pricing</button>
                </>
            </div>
            <div className={`rounded-lg bg-gray-50 px-6 pb-4 pt-8 flex flex-col gap-4 ${detailsPart ? "hidden" : "block"}`}>
                <>
                    <label htmlFor="offer type" className="text-lg font-semibold">
                        Offer Type
                    </label>
                    <div className="flex gap-8">
                        <div className="flex items-center">
                            <input
                                id="standard offer"
                                type="radio"
                                name="is_auction"
                                value="false"
                                checked={is_auction === false}
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
                                name="is_auction"
                                value="true"
                                checked={is_auction === true}
                                onChange={() => changeOfferType(true)}
                                className="h-4 w-4 text-blue-500 focus:ring-blue-400"
                            />
                            <label htmlFor="auction" className="ml-2 text-sm font-medium text-gray-900">
                                Auction
                            </label>
                        </div>
                    </div>
                    <div id="username-error" aria-live="polite" aria-atomic="true">
                        {state?.errors?.is_auction &&
                            state.errors.is_auction.map((error: string) => (
                            <p className="mt-2 text-sm text-red-500" key={error}>
                                {error}
                            </p>
                        ))}
                    </div>
                    <NumberInputField id="price" name="Price" />
                    <SelectionLabel id="margin" name="Margin ( % )" options={["8", "9", "10"]} required={!detailsPart}/>
                    <div className="flex items-center gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 text-blue-500" viewBox="0 0 20 20" fill="currentColor">
                            <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2h-1V9z" clipRule="evenodd" />
                        </svg>
                        <span className="text-sm text-gray-600">
                            This percentage will be taxed out of set price. Chosen margin will influence the priority of your offer in the search results. The higher the margin, the higher the priority.
                        </span>
                    </div>

                    {is_auction && (
                        <>
                            <DateSelectionField id="auction_end_date" name="Auction end date" />
                            <NumberInputField id="buy_now_auction_price" name="Buy now price [ optional ]" />
                        </>
                    )}

                    <button type="button" className="bg-blue-600 text-white rounded p-2" onClick={() => setDetailsPart(true)}>Back to Details</button>
                    <button type="submit" className="bg-blue-600 text-white rounded p-2">Confirm Offer</button>
                </>
            </div>
        </form>
    )
}