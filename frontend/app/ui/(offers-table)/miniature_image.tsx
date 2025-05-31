interface MiniatureImageProps {
    url: string;
}

export default function MainImage({ url }: MiniatureImageProps) {
    return (
        <img src={url} alt="Car image" className="h-35 w-70 object-cover rounded" />
    );
}