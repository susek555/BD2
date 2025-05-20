import { OfferDetails } from "@/app/lib/definitions/reviews";

export default function OfferDetailsTable( details: { details: OfferDetails[] }) {
    return (
        <div className="flex flex-col gap-4 w-full">
            {details.details.map((detail, index) => (
                <div key={index} className="bg-gray-200 p-2 rounded">
                    <div className="flex justify-between text-lg px-3">
                        <strong className="max-w-1/2 truncate">{detail.name}:</strong>
                        <span className="max-w-1/2 truncate text-right">{detail.value}</span>
                    </div>
                </div>
            ))}
        </div>
    );
}