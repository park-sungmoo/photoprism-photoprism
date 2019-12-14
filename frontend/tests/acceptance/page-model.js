import {Selector, t} from 'testcafe';

export default class Page {
    constructor() {
        this.view = Selector('div.p-view-select', {timeout: 15000});
        this.camera = Selector('div.p-camera-select', {timeout: 15000});
        this.countries = Selector('div.p-countries-select', {timeout: 15000});
        this.time = Selector('div.p-time-select', {timeout: 15000});
        this.search1 = Selector('div.p-search-field input', {timeout: 15000});
    }

    async setFilter(filter, option) {
        let filterSelector = "";

        switch (filter) {
            case 'view':
                filterSelector = 'div.p-view-select';
                break;
            case 'camera':
                filterSelector = 'div.p-camera-select';
                break;
            case 'time':
                filterSelector = 'div.p-time-select';
                break;
            case 'countries':
                filterSelector = 'div.p-countries-select';
                break;
            default:
                throw "unknown filter";
        }

        await t
            .click(filterSelector, {timeout: 15000});

        if (option) {
            await t
                .click(Selector('div.menuable__content__active div.v-select-list a').withText(option), {timeout: 15000})
        } else {
            await t
                .click(Selector('div.menuable__content__active div.v-select-list a').nth(1), {timeout: 15000})
        }
    }

    async search(term) {
        await t
            .typeText(this.search1, term)
            .pressKey('enter')
    }

    async openNav() {
        if (await Selector('button.p-navigation-show').visible) {
            await t.click(Selector('button.p-navigation-show'));
        } else if (await Selector('div.p-navigation-expand').exists) {
            await t.click(Selector('div.p-navigation-expand i'));
        }
    }

    async selectPhoto(nPhoto) {
        await t
        .hover(Selector('div[class="v-image__image v-image__image--cover"]', {timeout:4000}).nth(nPhoto))
        .click(Selector('.t-select.t-off'));
    }

    async unselectPhoto(nPhoto) {
        await t
            .hover(Selector('div[class="v-image__image v-image__image--cover"]', {timeout:4000}).nth(nPhoto))
            .click(Selector('.t-select.t-on'));
    }

    async likePhoto(nPhoto) {
        await t
            .hover(Selector('div[class="v-image__image v-image__image--cover"]', {timeout:4000}).nth(nPhoto))
            .click(Selector('.t-like.t-off'));
    }

    async login(password) {
        await t
            .typeText(Selector('input[type="password"]'), password)
            .pressKey('enter');
    }

    async logout() {
        await t
            .click(Selector('div.p-navigation-logout'));
    }
}
