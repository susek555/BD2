"use client";

import { BaseFilterTemplate } from "@/app/ui/(home)/base-filter-template/base-filter-template";
import { BaseRangeTemplate } from "./base-filter-template/base-range-template";
import { fetchFilterFields, prepareRangeFields } from "@/app/lib/data";
import { useEffect, useState } from "react";
import { FilterFieldData, RangeFieldData } from "@/app/lib/definitions";
import { usePathname, useSearchParams, useRouter } from "next/navigation";

export default function Filters() {
    const [filters, setFilters] = useState<FilterFieldData[]>([]);
    const [ranges, setRanges] = useState<RangeFieldData[]>([]);

    useEffect(() => {
        async function fetchData() {
            const data = await fetchFilterFields();
            setFilters(data);
            setRanges(prepareRangeFields());
        }
        fetchData();
    }, []);

    const searchParams = useSearchParams();
    const pathname = usePathname();
    const { replace } = useRouter();

    function handleFilterChange(name: string, selected: string[]) {
        const params = new URLSearchParams(searchParams);

        params.set('page', '1'); // Reset to the first page
        params.delete('query'); // Reset the search query
        if (selected.length > 0) {
            params.set(name, selected.join(","));
        } else {
            params.delete(name);
        }

        replace(`${pathname}?${params.toString()}`);
    }

    function handleRangeChange(name: string, range: { min: number; max: number }) {
        const params = new URLSearchParams(searchParams);

        params.set('page', '1'); // Reset to the first page
        params.delete('query'); // Reset the search query
        if (range.min !== 0) {
            params.set(`${name}_min`, range.min.toString());
        } else {
            params.delete(`${name}_min`);
        }

        if (range.max !== 0) {
            params.set(`${name}_max`, range.max.toString());
        } else {
            params.delete(`${name}_max`);
        }

        replace(`${pathname}?${params.toString()}`);
    }

    return (
        <>
            <p className="px-2">Filters</p>
            {filters.map((filter, index) => (
                <BaseFilterTemplate
                    key={index}
                    name={filter.fieldName}
                    options={filter.options}
                    onChange={handleFilterChange}
                />
            ))}
            {ranges.map((range, index) => (
                <BaseRangeTemplate
                    key={index}
                    fieldName={range.fieldName}
                    onChange={handleRangeChange}
                />
            ))}
        </>
    );
}