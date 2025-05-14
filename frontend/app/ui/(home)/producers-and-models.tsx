'use client'

import { syncFiltersWithParams } from "@/app/lib/(home)/syncWithParams";
import { ModelFieldData, FilterFieldData } from "@/app/lib/definitions";
import { BaseFilterTemplate } from "@/app/ui/(home)/base-filter-template/base-filter-template";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useState } from "react";

export default function ProducersAndModels({ producersAndModels }: { producersAndModels: ModelFieldData }) {
    const searchParams = useSearchParams();
    const pathname = usePathname();
    const { replace } = useRouter();

    const producers = syncFiltersWithParams([producersAndModels.producers], searchParams)[0];
    let models: FilterFieldData = { fieldName: "Models", options: [] };

    // Stan widoczności pola Models
    const [showModels, setShowModels] = useState<boolean>(!!producers.selected && producers.selected.length > 0);

    function handleChange(name: string, selected: string[]) {
        const params = new URLSearchParams(searchParams);
        const sanitizedName = name.replace(/\s+/g, ''); // Remove whitespaces from name

        params.set('page', '1'); // Reset to the first page
        if (selected.length > 0) {
            params.set(sanitizedName, selected.join(","));
            // Natychmiastowa aktualizacja widoczności
            if (name === producers.fieldName) {
                setShowModels(true);
            }
        } else {
            params.delete(sanitizedName);
            // Ukryj pole Models, jeśli nie ma wybranego producenta
            if (name === producers.fieldName) {
                setShowModels(false);
            }
        }

        console.log("producers.selected", producers.selected);

        replace(`${pathname}?${params.toString()}`);
    }

    return (
        <>
            <p className="px-2">Filters:</p>
            <BaseFilterTemplate
                name={producers.fieldName}
                options={producers.options}
                selected={producers.selected}
                onChange={handleChange}
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
