import {defineConfig} from 'sponsorkit';

const helpers = {
    avatar: {
        size: 45
    },
    boxWidth: 55,
    boxHeight: 55,
    container: {
        sidePadding: 30
    },
};

const coffee = {
    avatar: {
        size: 50
    },
    boxWidth: 65,
    boxHeight: 65,
    container: {
        sidePadding: 30
    },
};

const breakfast = {
    avatar: {
        size: 55
    },
    boxWidth: 75,
    boxHeight: 75,
    container: {
        sidePadding: 20
    },
    name: {
        maxLength: 10
    }
};

const costs = {
    avatar: {
        size: 65
    },
    boxWidth: 90,
    boxHeight: 80,
    container: {
        sidePadding: 30
    },
    name: {
        maxLength: 10
    }
};

const bronze = {
    avatar: {
        size: 85
    },
    boxWidth: 110,
    boxHeight: 100,
    container: {
        sidePadding: 30
    },
    name: {
        maxLength: 20
    }
};

const silver = {
    avatar: {
        size: 100
    },
    boxWidth: 110,
    boxHeight: 110,
    container: {
        sidePadding: 20
    },
    name: {
        maxLength: 20
    }
};

const gold = {
    avatar: {
        size: 150
    },
    boxWidth: 175,
    boxHeight: 175,
    container: {
        sidePadding: 25
    },
    name: {
        maxLength: 25
    }
};

const champion = {
    avatar: {
        size: 175
    },
    boxWidth: 200,
    boxHeight: 200,
    container: {
        sidePadding: 30
    },
    name: {
        maxLength: 30
    }
};

const partner = {
    avatar: {
        size: 200
    },
    boxWidth: 225,
    boxHeight: 225,
    container: {
        sidePadding: 40
    },
    name: {
        maxLength: 40
    },

};

export default defineConfig({
    github: {
        login: 'leaanthony',
        type: 'user',
    },

    // Rendering configs
    width: 800,
    formats: ['svg'],
    tiers: [
        {
            title: 'Helpers',
            preset: helpers,
            composeAfter: function (composer, tierSponsors, config) {
                composer.addSpan(20);
            }
        },
        {
            title: 'Buying Coffee',
            monthlyDollars: 5,
            preset: coffee,
            composeAfter: function (composer, tierSponsors, config) {
                composer.addSpan(20);
            }
        },
        {
            title: 'Buying Breakfast',
            monthlyDollars: 10,
            preset: breakfast,
            composeAfter: function (composer, tierSponsors, config) {
                composer.addSpan(20);
            }
        },
        {
            title: 'Covering Costs',
            monthlyDollars: 20,
            preset: costs,
            composeAfter: function (composer, tierSponsors, config) {
                composer.addSpan(20);
            }
        },
        {
            title: 'Bronze Sponsors',
            monthlyDollars: 50,
            preset: bronze,
            composeAfter: function (composer, tierSponsors, config) {
                composer.addSpan(20);
            }
        },
        {
            title: 'Silver Sponsors',
            monthlyDollars: 100,
            preset: silver,
            composeAfter: function (composer, tierSponsors, config) {
                composer.addSpan(20);
            }
        },
        {
            title: 'Gold Sponsors',
            monthlyDollars: 200,
            preset: gold,
            composeAfter: function (composer, tierSponsors, config) {
                composer.addSpan(20);
            }
        },
        {
            title: 'Champion',
            monthlyDollars: 500,
            preset: champion,
            composeAfter: function (composer, tierSponsors, config) {
                composer.addSpan(20);
            }
        },
        {
            title: 'Partner',
            monthlyDollars: 1000,
            preset: partner,
            composeAfter: function (composer, tierSponsors, config) {
                composer.addSpan(20);
            }
        },
    ],
});