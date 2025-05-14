'use client'

import { syncFiltersWithParams } from "@/app/lib/(home)/syncWithParams";
import { ModelFieldData, FilterFieldData } from "@/app/lib/definitions";
import { BaseFilterTemplate } from "@/app/ui/(home)/base-filter-template/base-filter-template";
import { get } from "http";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useState } from "react";

export default function ProducersAndModels({ producersAndModels }: { producersAndModels: ModelFieldData }) {
    const searchParams = useSearchParams();
    const pathname = usePathname();
    const { replace } = useRouter();

    function getAvailableModels(producersData: FilterFieldData, models: string[][]) : string[] {
        const selectedIndexes = (producersData.selected ?? [])
            .map(selectedItem => producersData.options.indexOf(selectedItem))
            .filter(index => index !== -1);

        const result = selectedIndexes.flatMap(index =>
            models[index].map(model => `${producersData.options[index]} ${model}`)
        );

        return result;
    }

    const producers = syncFiltersWithParams([producersAndModels.producers], searchParams)[0] as { fieldName: string; options: string[]; selected: string[] };
    const [models, setModels] = useState<FilterFieldData>(syncFiltersWithParams([{
        fieldName: "Models",
        options: getAvailableModels(producers, producersAndModels.models),
    }], searchParams)[0]);

    const [showModels, setShowModels] = useState<boolean>(!!producers.selected && producers.selected.length > 0);

    function handleProducersChange(name: string, selected: string[]) {
        handleChange(name, selected);

        if (selected.length > 0) {
            const updatedModels = getAvailableModels({ ...producers, selected }, producersAndModels.models);
            setModels({
                fieldName: "Models",
                options: updatedModels,
                selected: []
            });
            setShowModels(true);
        } else {
            setModels({ fieldName: "Models", options: [], selected: [] });
            setShowModels(false);
        }
    }

    function handleChange(name: string, selected: string[]) {
        const params = new URLSearchParams(searchParams);
        const sanitizedName = name.replace(/\s+/g, '');

        params.set('page', '1');
        if (selected.length > 0) {
            params.set(sanitizedName, selected.join(","));
        } else {
            params.delete(sanitizedName);
        }

        replace(`${pathname}?${params.toString()}`);
    }

    return (
        <>
            <p className="px-2">Filters:</p>
            <BaseFilterTemplate
                name={producers.fieldName}
                options={producers.options}
                selected={producers.selected}
                onChange={handleProducersChange}
            />
            {showModels && (
                <BaseFilterTemplate
                    name={models.fieldName}
                    options={models.options}
                    selected={models.selected}
                    onChange={handleChange}
                />
            )}
        </>
    );
}
