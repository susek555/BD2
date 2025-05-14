'use client'

import { syncFiltersWithParams } from "@/app/lib/(home)/syncWithParams";
import { ModelFieldData } from "@/app/lib/definitions"
import { BaseFilterTemplate } from "@/app/ui/(home)/base-filter-template/base-filter-template";
import { usePathname, useRouter, useSearchParams } from "next/navigation";

export default function ProducersAndModels({ producersAndModels }: { producersAndModels: ModelFieldData }) {
    const searchParams = useSearchParams();
    const pathname = usePathname();
    const { replace } = useRouter();

    const producers = syncFiltersWithParams([producersAndModels.producers], searchParams)[0];

    function handleProducerChange(name: string, selected: string[]) {
        const params = new URLSearchParams(searchParams);
        const sanitizedName = name.replace(/\s+/g, ''); // Remove whitespaces from name

        params.set('page', '1'); // Reset to the first page
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
                onChange={handleProducerChange} />
        </>
    )
}