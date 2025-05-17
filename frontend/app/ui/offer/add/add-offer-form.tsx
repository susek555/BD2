"use client"

import { AddOfferFormData, AddOfferFormState } from "@/app/lib/definitions";
import React, { startTransition, useActionState } from "react";
import { getAvailableModels } from "../../(home)/producers-and-models";
import { addOffer } from "@/app/actions/add-offer";
import { useDebouncedCallback } from "use-debounce";

export default function AddOfferForm({ inputsData } : { inputsData : AddOfferFormData}) {
    const initialState: AddOfferFormState = {
        errors: {},
        values: {is_auction: false},
    };

    const [state, action] = useActionState(addOffer, initialState);
    const [liveState, setLiveState] = React.useState<AddOfferFormState>(state);

    // React.useEffect(() => {
    //     if (state !== liveState) {
    //         setLiveState(state);
    //         console.log("State updated:", state);
    //     }
    // }, [state]);

    const updateField = (name: string, value: any) => {
        setLiveState((prev) => ({
            ...prev,
            values: {
                ...prev.values,
                [name]: value,
            },
        }));
    };

    const debouncedUpdateField = useDebouncedCallback((name: string, val: any) => {
        updateField(name, val);
    }, 500);

    const [availableModels, setAvailableModels] = React.useState<string[]>([]);

    const handleProducerChange = (newProducer: string) => {
        const models = getAvailableModels(
            { fieldName: "Producers", options: inputsData.producers, selected: [newProducer] },
            inputsData.models
        );
        setAvailableModels(models);
        updateField("producer", newProducer);
    };

    const [description, setDescription] = React.useState<string>(liveState.values?.description?.toString() ?? "");

    const changeOfferType = (isAuction: boolean) => {
        updateField("is_auction", isAuction);
    };

    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        const formData = new FormData(event.currentTarget);

        formData.append("isAuction", liveState.values?.isAuction == null ? "false" : liveState.values.isAuction.toString());
        const registraction_day = formData.get("registration_date-day");
        const registraction_month = formData.get("registration_date-month");
        const registraction_year = formData.get("registration_date-year");
        if (registraction_day && registraction_month && registraction_year) {
            const date = `${registraction_year.toString().padStart(4, "0")}-${registraction_month.toString().padStart(2, "0")}-${registraction_day.toString().padStart(2, "0")}`;
            formData.set("registration_date", date);
            formData.delete("registration_date-day");
            formData.delete("registration_date-month");
            formData.delete("registration_date-year");
            updateField("registration_date", date);
        }
        const auction_day = formData.get("auction_end_date-day");
        const auction_month = formData.get("auction_end_date-month");
        const auction_year = formData.get("auction_end_date-year");
        if (auction_day && auction_month && auction_year) {
            const auctionDate = `${auction_year.toString().padStart(4, "0")}-${auction_month.toString().padStart(2, "0")}-${auction_day.toString().padStart(2, "0")}`;
            formData.set("auction_end_date", auctionDate);
            formData.delete("auction_end_date-day");
            formData.delete("auction_end_date-month");
            formData.delete("auction_end_date-year");
            updateField("auction_end_date", auctionDate);
        }

        console.log("Old liveState:", liveState);

        startTransition(() => {
            action(formData)
        });

    }

    // UI

    function SelectionLabel({ id, name, options }: { id: string, name: string; options: string[] }) {
        return (
            <>
                <label htmlFor={id} className="text-lg font-semibold">
                    {name}
                </label>
                <select
                    id={id}
                    name={id}
                    className="border rounded p-2"
                    required
                    // defaultValue={
                    //     typeof state?.values?.[id] === "string"
                    //         ? state.values[id]
                    //         : ""
                    // }
                    value={liveState.values?.[id]?.toString() ?? ""}
                    {...(id === "producer" ? {
                        onChange: (e) => handleProducerChange(e.target.value)
                    } : {
                        onChange: (e) => updateField(id, e.target.value)
                    })}
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
        const [value, setValue] = React.useState<string>(liveState.values?.[id]?.toString() ?? "");

        function handleNumberChange(id: string, value: string) {
            const parsedValue = parseInt(value);
            if (!isNaN(parsedValue)) {
                debouncedUpdateField(id, parsedValue);
            } else {
                debouncedUpdateField(id, value);
            }
        }

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
                    value={value}
                    onChange={(e) => {
                        setValue(e.target.value);
                        handleNumberChange(id, e.target.value)
                    }}
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
        const [initialDate, setInitialDate] = React.useState<string>(liveState.values?.[id]?.toString() ?? "");

        const [year, setYear] = React.useState<string>(() => {
            return initialDate ? initialDate.split("-")[0] ?? "" : "";
        });
        const [month, setMonth] = React.useState<string>(() => {
            return initialDate ? initialDate.split("-")[1] ?? "" : "";
        });
        const [day, setDay] = React.useState<string>(() => {
            return initialDate ? initialDate.split("-")[2] ?? "" : "";
        });

        const prevInitialDateRef = React.useRef<string>("");

        React.useEffect(() => {
        if (liveState.values?.[id]?.toString() !== prevInitialDateRef.current) {
            const newInitialDate = liveState.values?.[id]?.toString() ?? "";
            setInitialDate(newInitialDate);
            const parts = newInitialDate.split("-");
            setYear(parts[0] ?? "");
            setMonth(parts[1] ?? "");
            setDay(parts[2] ?? "");
            prevInitialDateRef.current = newInitialDate;
        }
        }, [liveState.values, id]);

        React.useEffect(() => {
        if (
            year.length === 4 &&
            month.length > 0 && parseInt(month) >= 1 && parseInt(month) <= 12 &&
            day.length > 0 && parseInt(day) >= 1 && parseInt(day) <= 31
        ) {
            const newDate = `${year.padStart(4, "0")}-${month.padStart(2, "0")}-${day.padStart(2, "0")}`;
            if (newDate !== prevInitialDateRef.current) {
                debouncedUpdateField(id, newDate);
                setInitialDate(newDate);
                prevInitialDateRef.current = newDate;
            }
        }
        }, [year, month, day, id]);

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
                        value={day}
                        onChange={(e) => setDay(e.target.value)}
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
                        value={month}
                        onChange={(e) => setMonth(e.target.value)}
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
                        value={year}
                        onChange={(e) => setYear(e.target.value)}
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
        const [value, setValue] = React.useState<string>(liveState.values?.[id]?.toString() ?? "");

        function handleTextChange(id: string, value: string) {
            setValue(value);
            debouncedUpdateField(id, value);
        }

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
                    value={value}
                    onChange={(e) => handleTextChange(id, e.target.value)}
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
        <form className=" w-full md:w-200" onSubmit={handleSubmit}>
            <div className="rounded-lg bg-gray-50 px-6 pb-4 pt-8 flex flex-col gap-4">
                <label htmlFor="producer" className="text-lg font-semibold">Producer</label>
                <select
                    id="producer"
                    name="producer"
                    className="border rounded p-2"
                    required
                    value={liveState.values?.producer || ""}
                    onChange={e => handleProducerChange(e.target.value)}
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
                <SelectionLabel id="country" name="Country" options={inputsData.countries} />
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
                <TextInputField id="location" name="Location" />
                <TextInputField id="vin" name="VIN" />

                <label htmlFor="description" className="text-lg font-semibold">
                    Description
                </label>
                <textarea
                    id="description"
                    name="description"
                    className="border rounded p-2 h-32"
                    value={description}
                    onChange={(e) => {
                        debouncedUpdateField("description", e.target.value)
                        setDescription(e.target.value);}}
                    placeholder="Enter description..."
                    required
                ></textarea>

                <label htmlFor="images" className="text-lg font-semibold">
                    Images
                </label>
                TODO
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
                            checked={liveState.values?.is_auction === false}
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
                            checked={liveState.values?.is_auction === true}
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

                {liveState.values?.is_auction && (
                    <>
                        <DateSelectionField id="auction_end_date" name="Auction end date" />
                        <NumberInputField id="buy_now_auction_price" name="Buy now price" />
                    </>
                )}

                <button type="submit" className="bg-blue-600 text-white rounded p-2">Add offer</button>
            </div>
        </form>
    )
}