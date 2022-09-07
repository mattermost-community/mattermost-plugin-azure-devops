import React from 'react';

import './styles.scss';

type CardBody = {
    sectionHeading: string | JSX.Element;
    data: {
        icon: JSX.Element;
        label: string | JSX.Element;
    }[];
}

type CardProps = {
    cardHeader: string;
    cardSubHeader?: string;
    cardBody: CardBody[];
    onDelete?: (e: React.MouseEvent<HTMLButtonElement>) => void;
    onEdit?: (e: React.MouseEvent<HTMLButtonElement>) => void;
}

const Card = ({cardHeader, cardSubHeader, cardBody, onDelete, onEdit}: CardProps) => (
    <div className='details-card'>
        <div className='details-card__header d-flex align-items-center justify-content-between'>
            <h3 className='details-card__header-text text-ellipsis'>
                {cardHeader}
                {cardSubHeader && <span className='details-card__sub-header'>{cardSubHeader}</span>}
            </h3>
            <div>
                {onEdit && (
                    <button
                        onClick={onEdit}
                        className='details-card__action-btn details-card__action-btn--edit'
                    >
                        <i className='fa fa-pencil'/>
                    </button>
                )}
                {onDelete && (
                    <button
                        onClick={onDelete}
                        className='details-card__action-btn'
                    >
                        <i className='fa fa-trash'/>
                    </button>
                )}
            </div>
        </div>
        <div className='card-body'>
            {cardBody.map((body) => (
                <div
                    key={body.sectionHeading as string}
                    className='card-body__section'
                >
                    <h3 className='card-body__section-heading text-ellipsis'>{body.sectionHeading}</h3>
                    <ul className='card-body__list'>
                        {
                            body.data.map((listItem) => (
                                <li
                                    key={listItem.label as string}
                                    className='body-item'
                                >
                                    <p className='body-item__text text-ellipsis'>
                                        <span className='body-item__icon'>{listItem.icon}</span>
                                        <span>{listItem.label}</span>
                                    </p>
                                </li>
                            ))
                        }
                    </ul>
                </div>
            )) }
        </div>
    </div>
);

export default Card;
