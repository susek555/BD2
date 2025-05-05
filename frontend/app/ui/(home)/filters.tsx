"use client";

import { BaseFilterTemplate } from "@/app/ui/(home)/base-filter-template/base-filter-template";

export default function Filters() {

    function temp(selected: string[]) {}

    return (
        <>
            <p>Filters</p>
            <BaseFilterTemplate name="Model" options={["A5", "B4"]} onChange={temp}></BaseFilterTemplate>
        </>
    );
}


// TODO: implement