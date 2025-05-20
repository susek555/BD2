import { SaleOfferDetails } from "@/app/lib/definitions/sale-offer-details";

// Offer page

export async function fetchOfferDetails(id: string) : Promise<SaleOfferDetails | null> {
    // TODO connect API

    await new Promise(resolve => setTimeout(resolve, 1000));

    let data: SaleOfferDetails | null = null;

    if (id === "2") {
        data = {
            name: "Volkswagen Golf",
            price: 15000,
            isAuction: true,
            isActive: true,
            auctionData: {
            endDate: new Date("2025-07-07T12:00:00Z"),
            currentBid: 12000
            },
            imagesURLs: [
            "http://localhost:8081/test_car_image_1.webp",
            "http://localhost:8081/test_car_image_2.webp",
            "http://localhost:808/test_car_image_3.webp",
            ], // temp mock
            description: "The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology, and a timeless design, making it a popular choice among drivers worldwide.The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology, and a timeless design, making it a popular choice among drivers worldwide.The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology, and a timeless design, making it a popular choice among drivers worldwide.The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology,The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology, and a timeless design, making it a popular choice among drivers worldwide.The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology, and a timeless design, making it a popular choice among drivers worldwide.The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology, and a timeless design, making it a popular choice among drivers worldwide. and a timeless design, making it a popular choice among drivers worldwide.",
            details: [
                {
                    name: "Engine",
                    value: "2.0 TDI"
                },
                {
                    name: "Power",
                    value: "150 HP"
                },
                {
                    name: "Color",
                    value: "Blue"
                }
            ],
            sellerName: "Zwinny Ambroży",
            is_favourite: true,
            can_delete: false,
            can_edit: false,
        };
    }

    return data;
}


