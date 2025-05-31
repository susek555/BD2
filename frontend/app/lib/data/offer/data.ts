import { SaleOfferDetails } from '@/app/lib/definitions/sale-offer-details';
import { getOfferData, getSellerName } from '../../api/offer/fetchOffer';

// Offer page

export async function fetchOfferDetails(
  id: string,
): Promise<SaleOfferDetails | null> {

  const fetchedData = await getOfferData(parseInt(id, 10));

  if (!fetchedData) {
    return null;
  } else {
    const data: SaleOfferDetails = {
      id: fetchedData.id,
      name: `${fetchedData.brand} ${fetchedData.model}`,
      price: fetchedData.price,
      isAuction: fetchedData.is_auction,
      // isActive: fetchedData.is_active,
      // TODO include handling it
      isActive: true, // Temporary fix until backend is ready
      imagesURLs: fetchedData.images_urls,
      description: fetchedData.description,
      details: [
      {
        name: 'Color',
        value: fetchedData.color,
      },
      {
        name: 'Production Year',
        value: fetchedData.production_year?.toString() || ''
      },
      {
        name: 'Mileage',
        value: `${fetchedData.mileage} km`
      },
      {
        name: 'Fuel Type',
        value: fetchedData.fuel_type
      },
      {
        name: 'Transmission',
        value: fetchedData.transmission
      },
      {
        name: 'Number of gears',
        value: fetchedData.number_of_gears?.toString() || ''
      },
      {
        name: 'Number of doors',
        value: fetchedData.number_of_doors?.toString() || ''
      },
      {
        name: 'Number of seats',
        value: fetchedData.number_of_seats?.toString() || ''
      },
      {
        name: 'Engine Power',
        value: `${fetchedData.engine_power} HP`
      },
      {
        name: 'Engine Capacity',
        value: `${fetchedData.engine_capacity} L`
      },
      {
        name: 'Drive',
        value: fetchedData.drive
      },
      {
        name: 'VIN',
        value: fetchedData.vin
      },
      {
        name: 'Date of first registration',
        value: fetchedData.registration_date
      },
      {
        name: 'Registration number',
        value: fetchedData.registration_number || 'N/A'
      },
      {
        name: 'Date of offer issue',
        value: fetchedData.date_of_issue
      }
      ],
      sellerName: await getSellerName(fetchedData.user_id),
      sellerId: fetchedData.user_id,
      is_favourite: fetchedData.is_liked,
      can_delete: fetchedData.can_modify,
      can_edit: fetchedData.can_modify,
      auctionData: fetchedData.is_auction ? {
      endDate: new Date(fetchedData.auction_end_date),
      currentBid: fetchedData.current_bid || 0,
      } : undefined
    }
    return data;
  }


  //     name: 'Volkswagen Golf',
  //     price: 15000,
  //     isAuction: true,
  //     isActive: true,
  //     auctionData: {
  //       endDate: new Date('2025-07-07T12:00:00Z'),
  //       currentBid: 12000,
  //     },
  //     imagesURLs: [
  //       'http://localhost:8081/test_car_image_1.webp',
  //       'http://localhost:8081/test_car_image_2.webp',
  //       'http://localhost:808/test_car_image_3.webp',
  //     ], // temp mock
  //     description:
  //       'The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology, and a timeless design, making it a popular choice among drivers worldwide.The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology, and a timeless design, making it a popular choice among drivers worldwide.The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology, and a timeless design, making it a popular choice among drivers worldwide.The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology,The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology, and a timeless design, making it a popular choice among drivers worldwide.The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology, and a timeless design, making it a popular choice among drivers worldwide.The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology, and a timeless design, making it a popular choice among drivers worldwide. and a timeless design, making it a popular choice among drivers worldwide.',
  //     details: [
  //       {
  //         name: 'Engine',
  //         value: '2.0 TDI',
  //       },
  //       {
  //         name: 'Power',
  //         value: '150 HP',
  //       },
  //       {
  //         name: 'Color',
  //         value: 'Blue',
  //       },
  //     ],
  //     sellerName: 'Zwinny Ambroży',
  //     sellerId: 1,
  //     is_favourite: true,
  //     can_delete: false,
  //     can_edit: false,
  //   };
  // }

  // return data;
}
