'use client'

import { syncFiltersWithParams } from "@/app/lib/data/(home)/syncWithParams";
import { ModelFieldData, FilterFieldData } from "@/app/lib/definitions/reviews";
import { BaseFilterTemplate } from "@/app/ui/(filters-sidebar)/base-filter-template/base-filter-template";
import { usePathname, useSearchParams } from "next/navigation";
import { useState } from "react";

export function getAvailableModels(producersData: FilterFieldData, models: string[][]) : string[] {
        const selectedIndexes = (producersData.selected ?? [])
            .map(selectedItem => producersData.options.indexOf(selectedItem))
            .filter(index => index !== -1);

        const result = selectedIndexes.flatMap(index =>
            models[index].map(model => `${producersData.options[index]} ${model}`)
        );

        return result;
    }

export default function ProducersAndModels({ producersAndModels }: { producersAndModels: ModelFieldData }) {

    const searchParams = useSearchParams();
    const pathname = usePathname();

    const producers = syncFiltersWithParams([producersAndModels.producers], searchParams)[0] as { fieldName: string; options: string[]; selected: string[] };
    const [models, setModels] = useState<FilterFieldData>(syncFiltersWithParams([{
        fieldName: "Models",
        options: getAvailableModels(producers, producersAndModels.models),
        selected: []
    }], searchParams)[0]);

    const [showModels, setShowModels] = useState<boolean>(!!producers.selected && producers.selected.length > 0);
    const [isOpen, setIsOpen] = useState(false);

    // models field  - has to be in this file to ensure syncronization when change in
    //                 producers field clears selected values here

    function ModelsFilter() {

        const toggleDropdown = () => setIsOpen(!isOpen);

        const handleCheckboxChange = (option: string) => {
            const selectedOptions = models.selected || [];

            const updatedSelectedOptions = selectedOptions.includes(option)
                ? selectedOptions.filter((item) => item !== option)
                : [...selectedOptions, option];

            setModels({ ...models, selected: updatedSelectedOptions });
            handleModelChange("Models", updatedSelectedOptions);
        };

        return (
            <>
                <div className="base-filter-template border border-black-300 rounded px-2 py-1">
                    <button
                    className="flex justify-between items-center w-full"
                    onClick={toggleDropdown}
                    >
                    <span>{"Models"}</span>
                    <span>{isOpen ? '▲' : '▼'}</span>
                    </button>
                    {isOpen && (
                    <div className="filter-options mt-2">
                        {models.options.map((option) => (
                        <div
                            key={option}
                            className="filter-option flex justify-between items-center border border-gray-300 rounded px-2 py-1 mb-1"
                            onClick={() => handleCheckboxChange(option)}
                        >
                            <span>{option}</span>
                            <input
                            type="checkbox"
                            className="w-6 h-6 accent-green-500"
                            checked={models.selected?.includes(option)}
                            onChange={(e) => e.stopPropagation()}
                            />
                        </div>
                        ))}
                    </div>
                    )}
                </div>
            </>
        );
    }

    // actual component

    function handleProducersChange(name: string, selected: string[]) {
        const updatedModels = selected.length > 0
            ? getAvailableModels({ ...producers, selected }, producersAndModels.models)
            : [];

        setModels({
            fieldName: "Models",
            options: updatedModels,
            selected: []
        });

        setShowModels(updatedModels.length > 0);

        const params = new URLSearchParams(searchParams);
        const sanitizedName = name.replace(/\s+/g, '');

        params.set('page', '1');

        if (selected.length > 0) {
            params.set(sanitizedName, selected.join(","));
        } else {
            params.delete(sanitizedName);
        }

        params.delete("Models");

        window.history.replaceState(null, '', `${pathname}?${params.toString()}`);
    }

    function handleModelChange(name: string, selected: string[]) {
        const params = new URLSearchParams(searchParams);
        const sanitizedName = name.replace(/\s+/g, '');

        params.set('page', '1');
        if (selected.length > 0) {
            params.set(sanitizedName, selected.join(","));
        } else {
            params.delete(sanitizedName);
        }

        window.history.replaceState(null, '', `${pathname}?${params.toString()}`);
    }

    return (
        <>
            <h2 className='mb-2 text-sm font-semibold text-gray-700'>Filters</h2>
            <BaseFilterTemplate
                name={producers.fieldName}
                options={producers.options}
                selected={producers.selected}
                onChange={handleProducersChange}
            />
            {showModels && (
                <ModelsFilter />
            )}
        </>
    );
}
