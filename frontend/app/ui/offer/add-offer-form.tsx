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
    apiAction = offerActionEnum.ADD_OFFER
} : {
    inputsData : OfferFormData
    initialValues?: Partial<OfferFormState['values']>,
    apiAction?: boolean
}) {
    // Validate that ID is provided when edit mode is active
    // if (apiAction === offerActionEnum.EDIT_OFFER && id === undefined) {
    //     throw new Error("ID is required when editing an offer");
    // }


    const initialState: OfferFormState = {
        errors: {},
        values: initialValues as OfferFormState['values'],
    };

    const offerWrapper = (state: OfferFormState, formData: FormData) => {
            const progressState = formData.get('progressState') ? parseInt(formData.get('progressState') as string) : parseInt(OfferFormEnum.initialState.toString());
            formData.delete('progressState');

            const { progressState: result, offerFormState } = parseOfferForm(formData, progressState);


            setProgressState(result);

            if (result !== OfferFormEnum.readyToApi) {
                return offerFormState;
            }

            // Otherwise call the API action with validated data
            if (apiAction === offerActionEnum.ADD_OFFER) {
                return addOffer(offerFormState);
            } else {
                return editOffer(offerFormState);
            }
    };

    const [state, action] = useActionState(offerWrapper, initialState);

    useEffect(() => {
        if (state.values?.manufacturer) {
            setProducer(state.values.manufacturer.toString());
            handleProducerChange(state.values.manufacturer.toString());
        }
        if (typeof state.values?.is_auction === "boolean") {
            setIsAuction(state.values.is_auction);
        }

        console.log("State errors:", state.errors);
        console.log("State values:", state.values);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    , [state]);

    const [progressState, setProgressState] = React.useState<number>(OfferFormEnum.initialState);
    const [producer, setProducer] = React.useState<string>(state.values?.manufacturer?.toString() ?? "");

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
        console.log("Called handleSubmit");

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
        // Process auction date only when is_auction is true
        const is_auction_value = formData.get("is_auction") === "true";
        if (is_auction_value) {
            const auction_day = formData.get("date_end-day");
            const auction_month = formData.get("date_end-month");
            const auction_year = formData.get("date_end-year");
            const auction_hour = formData.get("date_end-hour");
            const auction_minute = formData.get("date_end-minute");

            if (auction_day && auction_month && auction_year && auction_hour && auction_minute) {
                const auctionDate = `${auction_hour.toString().padStart(2, "0")}:${auction_minute.toString().padStart(2, "0")} ${auction_year.toString().padStart(4, "0")}-${auction_month.toString().padStart(2, "0")}-${auction_day.toString().padStart(2, "0")}`;

                formData.set("date_end", auctionDate);
                formData.delete("date_end-day");
                formData.delete("date_end-month");
                formData.delete("date_end-year");
                formData.delete("date_end-hour");
                formData.delete("date_end-minute");
            }
        }

        formData.append('progressState', progressState.toString());

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
                            : typeof state.values?.[id] === "number"
                                ? state.values[id].toString()
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

    function DateSelectionField({ id, name, hasHour = false }: { id: string, name: string, hasHour?: boolean }) {
        const initialDate = state.values?.[id]?.toString() ?? "";
        const dateTimePattern = /^(\d{2}):(\d{2}) (\d{4})-(\d{2})-(\d{2})$/; // Pattern for "11:45 2026-12-21"
        const dateOnlyPattern = /^(\d{4})-(\d{2})-(\d{2})$/; // Pattern for "2026-12-21"

        const [
            initialHour = "",
            initialMinute = "",
            initialYear = "",
            initialMonth = "",
            initialDay = ""
        ] = initialDate
            ? dateTimePattern.test(initialDate)
                ? initialDate.match(dateTimePattern)?.slice(1, 6) || []
                : dateOnlyPattern.test(initialDate)
                    ? [...Array(2).fill(""), ...initialDate.match(dateOnlyPattern)?.slice(1, 4) || []]
                    : []
            : [];

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
                    {hasHour && (
                        <>
                        <div className="mx-2 flex items-center self-center">at</div>
                            <input
                                type="number"
                                id={`${id}-hour`}
                                name={`${id}-hour`}
                                min={0}
                                max={23}
                                className="border rounded p-2 w-20"
                                defaultValue={initialHour}
                                required
                                placeholder="Hour"
                            />
                            <input
                                type="number"
                                id={`${id}-minute`}
                                name={`${id}-minute`}
                                min={0}
                                max={59}
                                className="border rounded p-2 w-20"
                                defaultValue={initialMinute}
                                required
                                placeholder="Min"
                            />
                        </>
                    )}
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
        <>
            <form className=" w-full md:w-200" action={handleSubmit}>
                {/* Details Part */}
                <div className={`rounded-lg bg-gray-50 px-6 pb-4 pt-8 flex flex-col gap-4 ${progressState === OfferFormEnum.initialState ? "block" : "hidden"}`}>
                    <>
                        <label htmlFor="manufacturer" className="text-lg font-semibold">Manufacturer</label>
                        <select
                            id="manufacturer"
                            name="manufacturer"
                            className="border rounded p-2"
                            required
                            value={producer}
                            onChange={e => {
                                handleProducerChange(e.target.value);
                                setProducer(e.target.value); // Aktualizuj lokalny stan
                            } }
                        >
                            <option value="" disabled>Select a producer...</option>
                            {inputsData.producers.map((item: string) => (
                                <option key={item} value={item}>{item}</option>
                            ))}
                        </select>
                        <div id="username-error" aria-live="polite" aria-atomic="true">
                            {state?.errors?.manufacturer &&
                                state.errors.manufacturer.map((error: string) => (
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

                        <button type="submit" className="bg-blue-600 text-white rounded p-2">Set Pricing</button>
                    </>
                </div>
                {/* Pricing Part */}
                <div className={`rounded-lg bg-gray-50 px-6 pb-4 pt-8 flex flex-col gap-4 ${progressState === OfferFormEnum.pricingPart ? "block" : "hidden"}`}>
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
                                    className="h-4 w-4 text-blue-500 focus:ring-blue-400" />
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
                                    className="h-4 w-4 text-blue-500 focus:ring-blue-400" />
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
                        <SelectionLabel id="margin" name="Margin ( % )" options={["3", "5", "10"]} required={progressState >= OfferFormEnum.pricingPart} />
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
                                <DateSelectionField id="date_end" name="Auction end date" hasHour={true} />
                                <NumberInputField id="buy_now_price" name="Buy now price [ optional ]" />
                            </>
                        )}

                        <button type="button" className="bg-blue-600 text-white rounded p-2" onClick={() => setProgressState(OfferFormEnum.initialState)}>Back to Details</button>
                        <button type="submit" className="bg-blue-600 text-white rounded p-2">Upload images</button>
                    </>
                </div>
                {/* Images Part */}
                <div className={`rounded-lg bg-gray-50 px-6 pb-4 pt-8 flex flex-col gap-4 ${progressState === OfferFormEnum.imagesPart ? "block" : "hidden"}`}>
                    <>

                        <label htmlFor="images" className="text-lg font-semibold">
                            Images
                        </label>
                        <input type="file" id="images" name="images" className="border rounded p-2" multiple required={progressState === OfferFormEnum.imagesPart} />
                        <div id="username-error" aria-live="polite" aria-atomic="true">
                            {state?.errors?.images &&
                                state.errors.images.map((error: string) => (
                                    <p className="mt-2 text-sm text-red-500" key={error}>
                                        {error}
                                    </p>
                                ))}
                        </div>
                        <div className="flex items-center gap-2 mt-2 mb-4 p-3 bg-yellow-50 border border-yellow-300 rounded-md">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 text-yellow-600" viewBox="0 0 20 20" fill="currentColor">
                                <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                            </svg>
                            <span className="text-sm text-yellow-700">
                                Warning: If you go back to pricing, you&apos;ll need to select your images again.
                            </span>
                        </div>

                        <button type="button" className="bg-blue-600 text-white rounded p-2" onClick={() => setProgressState(OfferFormEnum.pricingPart)}>Back to Pricing</button>
                        <button type="submit" className="bg-blue-600 text-white rounded p-2">Submit Offer</button>
                    </>
                </div>
            </form>
            {/* Upload Results */}
            <div className={`rounded-lg bg-gray-50 px-6 pb-4 pt-8 flex flex-col gap-4 ${progressState === OfferFormEnum.readyToApi ? "block" : "hidden"}`}>
                    <>
                        {/* Upload failed */}
                        {state.errors?.upload_offer ? (
                            <div className="text-center space-y-4">
                                <div className="flex justify-center">
                                    <svg xmlns="http://www.w3.org/2000/svg" className="h-16 w-16 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                    </svg>
                                </div>
                                <h2 className="text-2xl font-semibold text-gray-800">Error Adding Offer</h2>
                                <p className="text-gray-600 max-w-md mx-auto">
                                    There was an issue submitting your offer. Please try again.
                                </p>
                                <div id="error-details" className="text-red-500 mt-2">
                                    {state.errors.upload_offer.map((error: string) => (
                                        <p key={error}>{error}</p>
                                    ))}
                                </div>
                                <div className="mt-6">
                                    <button
                                        type="button"
                                        onClick={() => setProgressState(OfferFormEnum.imagesPart)}
                                        className="inline-block bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-6 rounded-lg transition-colors"
                                    >
                                        Try Again
                                    </button>
                                </div>
                            </div>
                        ) : (
                            // Upload successful
                            <div className="text-center space-y-4">
                                <div className="flex justify-center">
                                    <svg xmlns="http://www.w3.org/2000/svg" className="h-16 w-16 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                                    </svg>
                                </div>
                                <h2 className="text-2xl font-semibold text-gray-800">Offer Successfully Added</h2>
                                <p className="text-gray-600 max-w-md mx-auto">
                                    Your offer has been successfully submitted and is now available in your listings.
                                </p>
                                <div className="mt-6">
                                    <a
                                        href="/account/listings"
                                        className="inline-block bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-6 rounded-lg transition-colors"
                                    >
                                        View My Listings
                                    </a>
                                </div>
                            </div>
                        )}
                    </>
            </div>
        </>
    )
}